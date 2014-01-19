// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "reflect"
    "sort"
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
    monitor, err := self.monitors.Get(
        name, func(_ interface{}) (interface{}, error) {
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

func MonitorStruct(data interface{}, cb func(name string, val float64)) {
    val := reflect.ValueOf(data)
    for val.Type().Kind() == reflect.Ptr {
        val = val.Elem()
    }
    typ := val.Type()
    if typ.Kind() != reflect.Struct {
        handleError(errors.ProgrammerError.New("not given a struct"))
        return
    }
    f64_type := reflect.TypeOf(float64(0))
    vals := make(map[string]float64, typ.NumField())
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        if field.Type.ConvertibleTo(f64_type) {
            vals[field.Name] = val.Field(i).Convert(f64_type).Float()
        }
    }
    MonitorMap(vals, cb)
}

func MonitorMap(data map[string]float64, cb func(name string, val float64)) {
    names := make([]string, 0, len(data))
    for name := range data {
        names = append(names, name)
    }
    sort.Strings(names)
    for _, name := range names {
        cb(name, data[name])
    }
}
