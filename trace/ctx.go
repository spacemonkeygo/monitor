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
	"fmt"

	"code.google.com/p/go.net/context"
	"github.com/spacemonkeygo/monitor"
)

type ctxKey int

const (
	spanKey ctxKey = iota
)

type spanCtx struct {
	span *Span
	context.Context
}

func (s *spanCtx) Value(key interface{}) interface{} {
	if key == spanKey {
		return s
	}
	return s.Context.Value(key)
}

func (s *spanCtx) String() string {
	if s.span.TraceDisabled() {
		return fmt.Sprintf("%v.WithDisabledSpan()", s.Context)
	}
	return fmt.Sprintf("%v.WithSpan(%#v)", s.Context, s.span.Export())
}

func getSpan(ctx context.Context) (s *Span, ctx_to_wrap context.Context) {
	if s, ok := ctx.(*spanCtx); ok && s != nil {
		return s.span, s.Context
	}
	if s, ok := ctx.Value(spanKey).(*spanCtx); ok && s != nil {
		return s.span, ctx
	}
	return nil, ctx
}

// TraceWithSpanName is like Trace, except you get to pick the Span name.
func (m *SpanManager) TraceWithSpanName(
	ctx *context.Context, name string) func(*error) {
	parent, parent_ctx := getSpan(*ctx)
	if parent == nil {
		s := m.NewTrace(name)
		new_ctx := &spanCtx{
			span:    s,
			Context: *ctx}
		*ctx = new_ctx
		return s.Observe()
	}
	if parent.TraceDisabled() {
		return func(*error) {}
	}
	s := parent.NewSpan(name)
	*ctx = &spanCtx{
		span:    s,
		Context: parent_ctx}
	return s.Observe()
}

// Trace is a utility that allows you to automatically create a Span (or a
// brand new trace, if needed) that observes the current function scope, given
// a function call context. The name of the Span is pulled from the current
// function name. See the example for usage.
func (m *SpanManager) Trace(ctx *context.Context) func(*error) {
	return m.TraceWithSpanName(ctx, monitor.CallerName())
}

// ContextWithSpan creates a new Context with the provided Span set as the
// current Span.
func ContextWithSpan(ctx context.Context, s *Span) context.Context {
	if s == nil {
		return ctx
	}
	return &spanCtx{
		span:    s,
		Context: ctx}
}

// SpanFromContext loads the current span from the current Context if one
// exists
func SpanFromContext(ctx context.Context) (s *Span, ok bool) {
	s, _ = getSpan(ctx)
	return s, s != nil
}
