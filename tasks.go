// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "flag"
    "sync"
    "time"

    "code.spacemonkey.com/go/errors"
)

var (
    maxErrorLength = flag.Int("monitor.max_error_length", 40,
        "the max length for an error name")
)

type TaskMonitor struct {
    mtx             sync.Mutex
    current         uint64
    highwater       uint64
    total_started   uint64
    total_completed uint64
    success         uint64
    timing          *ValueMonitor
    errors          map[string]uint64
}

func NewTaskMonitor() *TaskMonitor {
    return &TaskMonitor{
        errors: make(map[string]uint64),
        timing: NewValueMonitor()}
}

type TaskCtx struct {
    start   time.Time
    monitor *TaskMonitor
}

func (t *TaskMonitor) Start() *TaskCtx {
    t.mtx.Lock()
    t.current += 1
    t.total_started += 1
    if t.current > t.highwater {
        t.highwater = t.current
    }
    t.mtx.Unlock()
    return &TaskCtx{start: time.Now(), monitor: t}
}

func (c *TaskCtx) Finish(err error) {
    duration := time.Since(c.start)
    var error_name string
    if err != nil {
        if errors.HierarchicalError.Contains(err) {
            error_name = errors.GetClass(err).String()
        } else {
            error_name = err.Error()
        }
        if len(error_name) > *maxErrorLength {
            error_name = error_name[:*maxErrorLength]
        }
        error_name = SanitizeName(error_name)
    }

    c.monitor.mtx.Lock()
    c.monitor.current -= 1
    c.monitor.total_completed += 1
    if err != nil {
        c.monitor.errors[error_name] += 1
    } else {
        c.monitor.success += 1
    }
    c.monitor.mtx.Unlock()
    c.monitor.timing.Add(duration.Seconds())
}

func (self *MonitorGroup) Wrap() func(error) {
    return self.WrapNamed(CallerName())
}

func (self *MonitorGroup) WrapNamed(name string) func(error) {
    monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{}, error) {
        return NewTaskMonitor(), nil
    })
    if err != nil {
        handleError(err)
        return func(error) {}
    }
    task_monitor, ok := monitor.(*TaskMonitor)
    if !ok {
        handleError(errors.ProgrammerError.New(
            "monitor already exists with different type for name %s", name))
        return func(error) {}
    }
    ctx := task_monitor.Start()
    return ctx.Finish
}
