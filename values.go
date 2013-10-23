// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "sync"

    "code.spacemonkey.com/go/errors"
)

type ValueMonitor struct {
    mtx         sync.Mutex
    recent      float64
    count       uint64
    sum         float64
    sum_squared float64
    max         float64
    min         float64
}

func NewValueMonitor() *ValueMonitor {
    return &ValueMonitor{}
}

func (v *ValueMonitor) Add(val float64) {
    v.mtx.Lock()
    v.count += 1
    v.sum += val
    v.sum_squared += (val * val)
    v.recent = val
    if val > v.max {
        v.max = val
    }
    if val < v.min {
        v.min = val
    }
    v.mtx.Unlock()
}

func (self *MonitorGroup) Val(name string, val float64) {
    monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{}, error) {
        return NewValueMonitor(), nil
    })
    if err != nil {
        handleError(err)
        return
    }
    val_monitor, ok := monitor.(*ValueMonitor)
    if !ok {
        handleError(errors.ProgrammerError.New(
            "monitor already exists with different type for name %s", name))
        return
    }
    val_monitor.Add(val)
}
