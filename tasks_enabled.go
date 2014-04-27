// Copyright (C) 2013 Space Monkey, Inc.

// +build !no_mon

package monitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/SpaceMonkeyGo/errors"

	space_time "code.spacemonkey.com/go/space/time"
)

const (
	secondInMicroseconds     = 1000000
	microsecondInNanoseconds = 1000
)

type TaskCtx struct {
	start   time.Duration
	monitor *TaskMonitor
}

func (t *TaskMonitor) Start() func(*error) {
	ctx := t.NewContext()
	return func(e *error) { ctx.Finish(e, recover()) }
}

func (t *TaskMonitor) NewContext() *TaskCtx {
	t.mtx.Lock()
	t.current += 1
	t.total_started += 1
	if t.current > t.highwater {
		t.highwater = t.current
	}
	t.mtx.Unlock()
	return &TaskCtx{start: space_time.Monotonic(), monitor: t}
}

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

func (c *TaskCtx) Finish(err_ref *error, rec interface{}) {
	duration_nanoseconds := int64(space_time.Monotonic() - c.start)
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
	// we have to repanic
	if rec != nil {
		panic(rec)
	}
}
