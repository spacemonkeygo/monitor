// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "sync"

    "code.spacemonkey.com/go/errors"
)

type EventMonitor struct {
    mtx   sync.Mutex
    count uint64
}

func NewEventMonitor() *EventMonitor {
    return &EventMonitor{}
}

func (e *EventMonitor) Add() {
    e.mtx.Lock()
    e.count += 1
    e.mtx.Unlock()
}

func (e *EventMonitor) Stats(cb func(name string, val float64)) {
    e.mtx.Lock()
    count := e.count
    e.mtx.Unlock()
    cb("count", float64(count))
}

func (self *MonitorGroup) Event(name string) {
    monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{}, error) {
        return NewEventMonitor(), nil
    })
    if err != nil {
        handleError(err)
        return
    }
    event_monitor, ok := monitor.(*EventMonitor)
    if !ok {
        handleError(errors.ProgrammerError.New(
            "monitor already exists with different type for name %s", name))
        return
    }
    event_monitor.Add()
}
