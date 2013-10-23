// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "sync"

    "code.spacemonkey.com/go/errors"
)

type ChainedMonitor struct {
    mtx   sync.Mutex
    other Monitor
}

func NewChainedMonitor() *ChainedMonitor {
    return &ChainedMonitor{}
}

func (c *ChainedMonitor) Set(other Monitor) {
    c.mtx.Lock()
    c.other = other
    c.mtx.Unlock()
}

func (c *ChainedMonitor) Stats(cb func(name string, val float64)) {
    c.mtx.Lock()
    other := c.other
    c.mtx.Unlock()
    if other != nil {
        other.Stats(cb)
    }
}

func (self *MonitorGroup) Chain(name string, other Monitor) {
    name = SanitizeName(name)
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
    chain_monitor.Set(other)
}
