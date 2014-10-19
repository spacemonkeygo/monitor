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
	"io"
	"net/http"
	"sync"

	"code.google.com/p/go.net/context"
	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/monitor/trace/gen-go/zipkin"
)

// client stuff -----

// Client is an interface that matches an http.Client
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// TraceRequest will perform an HTTP request, creating a new Span for the HTTP
// request and sending the Span in the HTTP request headers.
// Compare to http.Client.Do.
func (m *SpanManager) TraceRequest(ctx context.Context, cl Client,
	req *http.Request, service *zipkin.Endpoint) (
	resp *http.Response, err error) {
	name := fmt.Sprintf("%s %s", req.Method, req.URL.String())
	s, ok := SpanFromContext(ctx)
	if ok {
		s = s.NewSpan(name)
	} else {
		s = m.NewTrace(name)
	}
	complete := s.ObserveService(service)
	s.Annotate("http.uri", req.URL.String(), service)
	s.Request().SetHeader(req.Header)
	resp, err = func() (resp *http.Response, err error) {
		defer errors.CatchPanic(&err)
		return cl.Do(req)
	}()
	if err != nil {
		complete(&err)
		return resp, err
	}
	s.Annotate("http.responsecode", int64(resp.StatusCode), service)
	current_body := resp.Body
	resp.Body = &wrappedBody{
		body:  current_body,
		close: func() { complete(nil) }}
	return resp, nil
}

type wrappedBody struct {
	body  io.ReadCloser
	close func()
	o     sync.Once
}

func (w *wrappedBody) Close() (err error) {
	err = w.body.Close()
	w.o.Do(w.close)
	return err
}

func (w *wrappedBody) Read(p []byte) (n int, err error) {
	n, err = w.body.Read(p)
	if err != nil {
		w.o.Do(w.close)
	}
	return n, err
}

// server stuff -----

// TraceHandler wraps a ContextHTTPHandler with a Span pulled from incoming
// requests, possibly starting new Traces if necessary.
func (m *SpanManager) TraceHandler(c ContextHTTPHandler) ContextHTTPHandler {
	return ContextHTTPHandlerFunc(func(
		ctx context.Context, w http.ResponseWriter, r *http.Request) {
		name := fmt.Sprintf("%s %s", r.Method, r.RequestURI)
		s := m.NewSpanFromRequest(name, RequestFromHeader(r.Header))
		defer s.Observe()(nil)
		s.Annotate("http.uri", r.RequestURI, nil)
		wrapped := &responseWriterObserver{w: w}
		c.ServeHTTP(ContextWithSpan(ctx, s), wrapped, r)
		s.Annotate("http.responsecode", fmt.Sprint(wrapped.StatusCode()), nil)
	})
}

type responseWriterObserver struct {
	w  http.ResponseWriter
	sc *int
}

func (w *responseWriterObserver) WriteHeader(status_code int) {
	w.sc = &status_code
	w.w.WriteHeader(status_code)
}

func (w *responseWriterObserver) Write(p []byte) (n int, err error) {
	if w.sc == nil {
		sc := 200
		w.sc = &sc
	}
	return w.w.Write(p)
}

func (w *responseWriterObserver) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriterObserver) StatusCode() int {
	if w.sc == nil {
		return 200
	}
	return *w.sc
}

// ContextHTTPHandler is like http.Handler, but expects a Context object
// as the first parameter.
type ContextHTTPHandler interface {
	ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

// ContextHTTPHandlerFunc is like http.HandlerFunc but for ContextHTTPHandlers
type ContextHTTPHandlerFunc func(
	ctx context.Context, w http.ResponseWriter, r *http.Request)

func (f ContextHTTPHandlerFunc) ServeHTTP(ctx context.Context,
	w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

// ContextWrapper will turn a ContextHTTPHandler into an http.Handler by
// passing a new Context into every request.
func ContextWrapper(h ContextHTTPHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(context.Background(), w, r)
	})
}
