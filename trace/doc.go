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

/*
Package trace is a distributed tracing and Zipkin client library for Go.

Background

For background, Zipkin (http://twitter.github.io/zipkin/) is Twitter's
implementation of the Google Dapper paper
(http://research.google.com/pubs/pub36356.html) for collecting tracing
information from a distributed system. Both Zipkin and Dapper instrument
RPC layers for collecting tracing information from the system as messages pass
between services. Please read more at the aforementioned websites.

This library is a Go client library to assist in integrating with Zipkin.

This library, for most common uses, relies heavily on Google's Context objects.
See http://blog.golang.org/context for more information there, but essentially
this library works best if you are already passing Context objects through most
of your callstacks.

Full example

See https://github.com/jtolds/go-zipkin-sample for a set of toy example
programs that use this library, or
https://raw.githubusercontent.com/jtolds/go-zipkin-sample/master/screenshot.png
for a screenshot of the Zipkin user interface after collecting a trace
from the sample application.

Basic usage

At a basic level, all you need to do to use this library to interface with
Zipkin are the two functions TraceHandler and TraceRequest.

Here's an example client:

  func MyOperation(ctx context.Context) error {
    req, err := http.NewRequest("GET", "http://my.url.tld/resource", nil)
    if err != nil {
      return err
    }
    resp, err := trace.TraceRequest(ctx, http.DefaultClient, req)
    if err != nil {
      return err
    }
    defer resp.Body.Close()

    // do stuff
  }

And an example server:

  func MyServer(addr string) error {
    return http.ListenAndServe(addr, trace.ContextWrapper(
        trace.TraceHandler(trace.ContextHTTPHandlerFunc(
        func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
          MyOperation(ctx)
          fmt.Fprintf(w, "hello\n")
        }))))
  }

It is important that the context objects are passed all the way through your
application logic from server to client to get the full effect.

In-process tracing

You may want to get tracing information out of operations and services within
a process, instead of only at RPC boundaries. In this scenario, there is one
more function used for creating new Spans within the same process.

For each function you want dedicated tracing information for, you can call
the Trace function like so:

  func MyTask(ctx context.Context) (result int, err error) {
    defer trace.Trace(&ctx)(&err)

    result, err = OtherTask1(ctx)
    if err != nil {
      return 0, err
    }

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
      defer wg.Done()
      OtherTask2(ctx)
    }()

    go func() {
      defer wg.Done()
      OtherTask3(ctx)
    }()

    wg.Wait()

    return result, nil
  }

Here, Trace modifies the current context to include a new Span named after the
calling function (MyTask). Your tracing collector will then receive this Span
and include annotations about when each sampled Span started, when it finished,
if it had a non-nil error or had a panic, what the error type was (if
github.com/spacemonkeygo/errors can identify it), and pass the Span along to
subcalls, for if they have their own spans.

If you don't like the automatic Span naming, you can use TraceWithSpanNamed
instead.

Process setup

Every process that sends Spans will need to be configured with Configure and
RegisterTraceCollector, so make sure to call those functions appropriately
early in your process lifetime.

Other

See https://github.com/itszero/docker-zipkin for easy Zipkin setup.

*/
package trace
