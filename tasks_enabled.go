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

// +build !no_mon

package monitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/monotime"
)

const (
	secondInMicroseconds     = 1000000
	microsecondInNanoseconds = 1000
)

// TaskCtx keeps track of a task as it is running.
type TaskCtx struct {
	start   time.Duration
	monitor *TaskMonitor
}

// Start is a helper method for watching a task in a less error-prone way.
// Managing a task context yourself is tricky to get right - recover only works
// in deferred methods. Call out of a method that was deferred and it no longer
// works! See the example.
func (t *TaskMonitor) Start() func(*error) {
	ctx := t.NewContext()
	return func(e *error) { ctx.Finish(e, recover()) }
}

// NewContext creates a new context that is watching a live task. See Start
// or MonitorGroup.Task
func (t *TaskMonitor) NewContext() *TaskCtx {
	t.mtx.Lock()
	t.current += 1
	t.total_started += 1
	if t.current > t.highwater {
		t.highwater = t.current
	}
	t.mtx.Unlock()
	return &TaskCtx{start: monotime.Monotonic(), monitor: t}
}

// Stats conforms to the Monitor interface
func (t *TaskMonitor) Stats(cb func(name string, val float64)) {
	t.mtx.Lock()
	current := t.current
	highwater := t.highwater
	total_started := t.total_started
	total_completed := t.total_completed
	success := t.success
	panics := t.panics
	error_counts := make(map[string]uint64, len(t.errors))
	for error, count := range t.errors {
		error_counts[error] = count
	}
	t.mtx.Unlock()

	errors := make([]string, 0, len(error_counts))
	for error := range error_counts {
		errors = append(errors, error)
	}
	sort.Strings(errors)

	cb("current", float64(current))
	for _, error := range errors {
		cb(fmt.Sprintf("error_%s", error), float64(error_counts[error]))
	}
	cb("highwater", float64(highwater))
	cb("panics", float64(panics))
	cb("success", float64(success))

	if len(errors) > 0 {
		t.error_timing.Stats(func(name string, val float64) {
			if name != "count" {
				// these values are in microseconds, convert to seconds
				cb(fmt.Sprintf("time_error_%s", name), val/secondInMicroseconds)
			}
		})
	}
	if success > 0 {
		t.success_timing.Stats(func(name string, val float64) {
			if name != "count" {
				// these values are in microseconds, convert to seconds
				cb(fmt.Sprintf("time_success_%s", name), val/secondInMicroseconds)
			}
		})
	}
	if total_completed > 0 {
		t.total_timing.Stats(func(name string, val float64) {
			if name != "count" {
				// these values are in microseconds, convert to seconds
				cb(fmt.Sprintf("time_total_%s", name), val/secondInMicroseconds)
			}
		})
	}
	cb("total_completed", float64(total_completed))
	cb("total_started", float64(total_started))
}

// Finish records a successful task completion. You must pass a pointer to
// the named error return value (or nil if there isn't one) and the result
// of recover() out of the method that was deferred for this to work right.
// Finish will re-panic any recovered panics (provided it wasn't a nil panic)
// after bookkeeping.
func (c *TaskCtx) Finish(err_ref *error, rec interface{}) {
	duration_nanoseconds := int64(monotime.Monotonic() - c.start)
	var error_name string
	var err error
	if err_ref != nil {
		err = *err_ref
	}
	if rec != nil {
		var ok bool
		err, ok = rec.(error)
		if !ok || err == nil {
			err = errors.PanicError.New("%v", rec)
		}
	}
	if err != nil {
		error_name = errors.GetClass(err).String()
		if len(error_name) > *maxErrorLength {
			error_name = error_name[:*maxErrorLength]
		}
		error_name = SanitizeName(error_name)
	}

	// we keep granularity on the order microseconds, which should keep
	// sum_squared useful
	duration_microseconds := int64(duration_nanoseconds /
		microsecondInNanoseconds)

	c.monitor.mtx.Lock()
	c.monitor.current -= 1
	c.monitor.total_completed += 1
	if err != nil {
		c.monitor.errors[error_name] += 1
		if rec != nil {
			c.monitor.panics += 1
		}
		c.monitor.error_timing.Add(duration_microseconds)
	} else {
		c.monitor.success_timing.Add(duration_microseconds)
		c.monitor.success += 1
	}
	c.monitor.mtx.Unlock()
	c.monitor.total_timing.Add(duration_microseconds)

	// doh, we didn't actually want to stop the panic codepath.
	// we have to repanic. Oh and great, panics can be nil. Welp!
	if rec != nil {
		panic(rec)
	}
}
