// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"reflect"
	"sort"
	"sync"

	"github.com/SpaceMonkeyGo/errors"
)

// ChainedMonitor is a monitor that simply wraps another monitor, while
// allowing for atomic monitor changing.
type ChainedMonitor struct {
	mtx   sync.Mutex
	other Monitor
}

// NewChainedMonitor returns a ChainedMonitor
func NewChainedMonitor() *ChainedMonitor {
	return &ChainedMonitor{}
}

// Set replaces the ChainedMonitor's existing monitor with other
func (c *ChainedMonitor) Set(other Monitor) {
	c.mtx.Lock()
	c.other = other
	c.mtx.Unlock()
}

// Stats conforms to the Monitor interface, and passes the call to the chained
// monitor.
func (c *ChainedMonitor) Stats(cb func(name string, val float64)) {
	c.mtx.Lock()
	other := c.other
	c.mtx.Unlock()
	if other != nil {
		other.Stats(cb)
	}
}

// MonitorStruct uses reflection to walk the structure data and call cb with
// every field and value on the struct that's castable to float64.
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

// MonitorMap sends a map of keys and values to a callback.
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
