// Copyright (C) 2013 Space Monkey, Inc.

// +build !no_mon

package client

import (
	"fmt"
	"strings"

	"code.spacemonkey.com/go/errors"
)

func (g *MonitorGroup) Stats(cb func(name string, val float64)) {
	snapshot := g.monitors.Snapshot()
	for _, name := range sortedStringKeys(snapshot) {
		cache_val := snapshot[name]
		mon, ok := cache_val.(Monitor)
		if !ok {
			continue
		}
		mon.Stats(func(subname string, val float64) {
			cb(fmt.Sprintf("%s.%s.%s", g.group_name, name, subname), val)
		})
	}
}

func (g *MonitorGroup) Datapoints(reset bool, cb func(name string,
	data [][]float64, total uint64, clipped bool, fraction float64)) {
	snapshot := g.collectors.Snapshot()
	for _, name := range sortedStringKeys(snapshot) {
		cache_val := snapshot[name]
		collector, ok := cache_val.(DataCollection)
		if !ok {
			continue
		}
		collector.Datapoints(reset, func(subname string, data [][]float64,
			total uint64, clipped bool, fraction float64) {
			cb(fmt.Sprintf("%s.%s.%s", g.group_name, name, subname), data,
				total, clipped, fraction)
		})
	}
}

func (self *MonitorGroup) Task() func(*error) {
	caller_name := CallerName()
	idx := strings.LastIndex(caller_name, "/")
	if idx >= 0 {
		caller_name = caller_name[idx+1:]
	}
	idx = strings.Index(caller_name, ".")
	if idx >= 0 {
		caller_name = caller_name[idx+1:]
	}
	return self.TaskNamed(caller_name)
}

func (self *MonitorGroup) TaskNamed(name string) func(*error) {
	name = SanitizeName(name)
	monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{}, error) {
		return NewTaskMonitor(), nil
	})
	if err != nil {
		handleError(err)
		return func(*error) {}
	}
	task_monitor, ok := monitor.(*TaskMonitor)
	if !ok {
		handleError(errors.ProgrammerError.New(
			"monitor already exists with different type for name %s", name))
		return func(*error) {}
	}
	return task_monitor.Start()
}

func (self *MonitorGroup) DataTask() func(*error) {
	// TODO: send data points
	caller_name := CallerName()
	idx := strings.LastIndex(caller_name, "/")
	if idx >= 0 {
		caller_name = caller_name[idx+1:]
	}
	idx = strings.Index(caller_name, ".")
	if idx >= 0 {
		caller_name = caller_name[idx+1:]
	}
	return self.TaskNamed(caller_name)
}

func (self *MonitorGroup) Data(name string, val ...float64) {
	name = SanitizeName(name)
	monitor, err := self.collectors.Get(name, func(_ interface{}) (interface{}, error) {
		return NewDatapointCollector(*collectionFraction, *collectionMax), nil
	})
	if err != nil {
		handleError(err)
		return
	}
	datapoint_collector, ok := monitor.(*DatapointCollector)
	if !ok {
		handleError(errors.ProgrammerError.New(
			"monitor already exists with different type for name %s", name))
		return
	}
	datapoint_collector.Add(val...)
}

func (self *MonitorGroup) Event(name string) {
	name = SanitizeName(name)
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

func (self *MonitorGroup) Val(name string, val float64) {
	name = SanitizeName(name)
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
