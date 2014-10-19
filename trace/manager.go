// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"sync"

	"code.google.com/p/go.net/context"
	"github.com/spacemonkeygo/monitor/trace/gen-go/zipkin"
)

// TraceCollector is an interface dealing with completed Spans on a
// SpanManager. See RegisterTraceCollector.
type TraceCollector interface {
	// Collect gets called with a Span whenever a Span is completed on a
	// SpanManager.
	Collect(span *zipkin.Span)
}

// TraceCollectorFunc is for closures that match the TraceCollector interface.
type TraceCollectorFunc func(span *zipkin.Span)

func (f TraceCollectorFunc) Collect(span *zipkin.Span) { f(span) }

// SpanManager creates and configures settings about Spans. Create one with
// NewSpanManager
type SpanManager struct {
	mtx                sync.Mutex
	default_local_host *zipkin.Endpoint
	trace_collectors   []TraceCollector
	trace_fraction     float64
	trace_debug        bool
}

// NewSpanManager creates a new SpanManager. No traces will be collected by
// default until Configure is called.
func NewSpanManager() *SpanManager { return &SpanManager{} }

// Configure configures a SpanManager. trace_fraction is the fraction of new
// traces that will be collected (between 0 and 1, inclusive). trace_debug
// is whether or not traces will have the debug flag set, which controls
// whether or not the trace collector will be allowed to sample them on its
// own. default_local_host is the annotation host endpoint to set when one
// isn't otherwise provided.
func (m *SpanManager) Configure(trace_fraction float64, trace_debug bool,
	default_local_host *zipkin.Endpoint) {
	m.mtx.Lock()
	m.trace_fraction = trace_fraction
	m.trace_debug = trace_debug
	m.default_local_host = default_local_host
	m.mtx.Unlock()
}

// RegisterTraceCollector takes a TraceCollector and calls Collect on it
// whenever a Span from this SpanManager is complete.
func (m *SpanManager) RegisterTraceCollector(collector TraceCollector) {
	m.mtx.Lock()
	m.trace_collectors = append(m.trace_collectors, collector)
	m.mtx.Unlock()
}

// NewSampledTrace creates a new span that begins a trace that is being sampled
// without consulting the configured trace_fraction. span_name names the first
// span of the trace, and debug controls whether or not the span collector is
// allowed to sample the trace on its own.
func (m *SpanManager) NewSampledTrace(span_name string, debug bool) *Span {
	trace_id := Rng.Int63() + 1
	return &Span{
		data: zipkin.Span{
			TraceId: trace_id,
			Name:    span_name,
			Id:      trace_id,
			Debug:   debug},
		manager: m}
}

// NewTrace creates a new span that begins a trace, after consulting the
// SpanManager's configured trace_fraction and trace_debug settings. The trace
// may or may not actually be sampled. span_name is the name of the beginning
// Span.
func (m *SpanManager) NewTrace(span_name string) *Span {
	m.mtx.Lock()
	trace_fraction := m.trace_fraction
	trace_debug := m.trace_debug
	m.mtx.Unlock()
	if Rng.Float64() >= trace_fraction {
		return NewDisabledTrace()
	}
	return m.NewSampledTrace(span_name, trace_debug)
}

// NewSpanFromRequest creates a new span, and possibly a new trace, given
// whatever was supplied in the incoming request.
func (m *SpanManager) NewSpanFromRequest(name string, req Request) *Span {
	if req.Sampled != nil && !*req.Sampled {
		return NewDisabledTrace()
	}

	if req.TraceId == nil || req.SpanId == nil {
		return m.NewTrace(name)
	}

	flags := int64(0)
	if req.Flags != nil {
		flags = *req.Flags
	}

	return &Span{
		data: zipkin.Span{
			TraceId:  *req.TraceId,
			Name:     name,
			Id:       *req.SpanId,
			ParentId: req.ParentId,
			Debug:    flags&1 > 0},
		server:  true,
		manager: m}
}

func (m *SpanManager) collect(s *Span) {
	m.mtx.Lock()
	collectors := m.trace_collectors
	m.mtx.Unlock()
	data := s.Export()
	for _, collector := range collectors {
		collector.Collect(data)
	}
}

func (m *SpanManager) defaultLocalHost() (rv *zipkin.Endpoint) {
	m.mtx.Lock()
	rv = m.default_local_host
	m.mtx.Unlock()
	return rv
}

var (
	DefaultManager = NewSpanManager()

	Configure              = DefaultManager.Configure
	NewSampledTrace        = DefaultManager.NewSampledTrace
	NewSpanFromRequest     = DefaultManager.NewSpanFromRequest
	NewTrace               = DefaultManager.NewTrace
	RegisterTraceCollector = DefaultManager.RegisterTraceCollector
	TraceHandler           = DefaultManager.TraceHandler
	TraceRequest           = DefaultManager.TraceRequest
	TraceWithSpanNamed     = DefaultManager.TraceWithSpanNamed
)

// Trace calls Trace on the DefaultManager
func Trace(ctx *context.Context) func(*error) {
	return DefaultManager.TraceWithSpanNamed(ctx, CallerName())
}
