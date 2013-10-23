// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "sync"

    "code.spacemonkey.com/go/errors"
)

type ChainedMonitor struct {
    mtx sync.Mutex
    cb  func() float64
}

func NewChainedMonitor() *ChainedMonitor {
    return &ChainedMonitor{}
}

func (c *ChainedMonitor) Set(cb func() float64) {
    c.mtx.Lock()
    c.cb = cb
    c.mtx.Unlock()
}

func (self *MonitorGroup) Chain(name string, val func() float64) {
    monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{}, error) {
        return NewChainedMonitor(), nil
    })
    if err != nil {
        handleError(err)
        return
    }
    chain_monitor, ok := monitor.(*ChainedMonitor)
    if !ok {
        handleError(errors.ProgrammerError.New(
            "monitor already exists with different type for name %s", name))
        return
    }
    chain_monitor.Set(val)
}
