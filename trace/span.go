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
	"time"

	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/monotime"
	"gopkg.in/spacemonkeygo/monitor.v1/trace/gen-go/zipkin"
)

// Span represents a given task or request within a full trace.
type Span struct {
	disabled bool
	// mtx just covers annotations, since everything else is read only
	mtx     sync.Mutex
	data    zipkin.Span
	server  bool
	manager *SpanManager
}

// Trace disabled returns whether the trace is even active. A disabled trace
// causes many operations to be a no-op.
func (s *Span) TraceDisabled() bool { return s.disabled }

// TraceId is the id of the given trace, if not disabled.
func (s *Span) TraceId() int64 { return s.data.TraceId }

// SpanId is the id of the given span, if not disabled.
func (s *Span) SpanId() int64 { return s.data.Id }

// Name is the name of the given span, if not disabled.
func (s *Span) Name() string { return s.data.Name }

// ParentId is the id of the parent span in the given trace, if not disabled.
func (s *Span) ParentId() *int64 { return s.data.ParentId }

// Debug is whether or not the trace collector is allowed to sample this trace
// on its own.
func (s *Span) Debug() bool { return s.data.Debug }

// AnnotateTimestamp annotates a given span with a timestamp. duration and
// host are optional.
func (s *Span) AnnotateTimestamp(key string, now time.Time,
	duration *time.Duration, host *zipkin.Endpoint) {
	if s.disabled {
		return
	}
	var duration_int32_ptr *int32
	if duration != nil {
		duration_int32 := int32(int64(*duration) / 1000)
		duration_int32_ptr = &duration_int32
	}
	if host == nil {
		host = s.manager.defaultLocalHost()
	}
	s.mtx.Lock()
	s.data.Annotations = append(s.data.Annotations, &zipkin.Annotation{
		Timestamp: now.UnixNano() / 1000,
		Value:     key,
		Duration:  duration_int32_ptr,
		Host:      host})
	s.mtx.Unlock()
}

// Annotate annotates a given span with an arbitrary value. host is optional.
// Annotate is a no-op unless val is of type time.Time, []byte, or string.
func (s *Span) Annotate(key string, val interface{}, host *zipkin.Endpoint) {
	if s.disabled {
		return
	}
	var serialized []byte
	var serialized_type zipkin.AnnotationType
	switch v := val.(type) {
	case time.Time:
		s.AnnotateTimestamp(key, v, nil, host)
		return
	case *time.Time:
		if v != nil {
			s.AnnotateTimestamp(key, *v, nil, host)
		}
		return
	case []byte:
		serialized_type = zipkin.AnnotationType_BYTES
		serialized = v
	case string:
		serialized_type = zipkin.AnnotationType_STRING
		serialized = []byte(v)
	default:
		return
	}
	if host == nil {
		host = s.manager.defaultLocalHost()
	}
	s.mtx.Lock()
	s.data.BinaryAnnotations = append(s.data.BinaryAnnotations,
		&zipkin.BinaryAnnotation{
			Key:            key,
			Value:          serialized,
			AnnotationType: serialized_type,
			Host:           host})
	s.mtx.Unlock()
}

// NewDisabledTrace creates a new Span that is disabled.
func NewDisabledTrace() *Span {
	return &Span{disabled: true}
}

// NewSpan creates a new Span off of the given parent Span.
func (parent *Span) NewSpan(name string) *Span {
	if parent.disabled {
		return parent
	}
	return &Span{
		data: zipkin.Span{
			TraceId:  parent.data.TraceId,
			Name:     name,
			Id:       Rng.Int63() + 1,
			ParentId: &parent.data.Id,
			Debug:    parent.data.Debug},
		manager: parent.manager}
}

// Export will take a Span and return a serializable thrift object.
func (s *Span) Export() *zipkin.Span {
	s.mtx.Lock()
	var copy zipkin.Span = s.data
	s.mtx.Unlock()
	return &copy
}

// Request will return a Request for RPC purposes based off of the existing
// Span
func (s *Span) Request() Request {
	sampled := !s.disabled
	flags := int64(0)
	if s.data.Debug {
		flags = 1
	}
	return Request{
		TraceId:  &s.data.TraceId,
		SpanId:   &s.data.Id,
		ParentId: s.data.ParentId,
		Sampled:  &sampled,
		Flags:    &flags}
}

// Observe is meant to watch a Span over a given Span duration.
func (s *Span) Observe() func(errptr *error) {
	return s.ObserveService(nil)
}

// ObserveService is like Observe, but uses a provided host instead of the
// SpanManager's default local host for annotations.
func (s *Span) ObserveService(service *zipkin.Endpoint) func(errptr *error) {
	if s.disabled {
		return func(*error) {}
	}
	if s.server {
		s.AnnotateTimestamp(zipkin.SERVER_RECV, monotime.Now(), nil, service)
	} else {
		s.AnnotateTimestamp(zipkin.CLIENT_SEND, monotime.Now(), nil, service)
	}
	return func(errptr *error) {
		end := monotime.Now()
		rec := recover()
		if rec != nil {
			s.AnnotateTimestamp("failed", end, nil, service)
			s.AnnotateTimestamp("panic", end, nil, service)
		} else {
			if errptr != nil && *errptr != nil {
				s.AnnotateTimestamp("failed", end, nil, service)
				s.Annotate("error", errors.GetClass(*errptr).String(), service)
			}
		}
		if s.server {
			s.AnnotateTimestamp(zipkin.SERVER_SEND, end, nil, service)
		} else {
			s.AnnotateTimestamp(zipkin.CLIENT_RECV, end, nil, service)
		}
		s.manager.collect(s)
		if rec != nil {
			panic(rec)
		}
	}
}
