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
	"strings"

	"github.com/spacemonkeygo/errors"
)

// Stats conforms to the Monitor interface. Stats aggregates all statistics
// attatched to this group.
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

// Datapoints conforms to the DataCollection interface. Datapoints aggregates
// all datasets attached to this group.
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

// Task allows you to monitor a specific function. Task automatically chooses
// a name for you based on the callstack and creates a TaskMonitor for you by
// that name if one doesn't already exist. If you'd like to pick your own
// metric name (and improve performance), use TaskNamed. Please see the
// example.
//
// N.B.: Error types are best tracked when you're using Space Monkey's
// hierarchical error package: http://github.com/spacemonkeygo/errors
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

// TaskNamed works just like Task without any automatic name selection
func (self *MonitorGroup) TaskNamed(name string) func(*error) {
	name = SanitizeName(name)
	monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{},
		error) {
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

// DataTask works just like Task, but automatically makes datapoints about
// the task in question. It's a hybrid of Data and Task.
func (self *MonitorGroup) DataTask() func(*error) {
	// TODO: actually send data points
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

// Data takes a name, makes a DataCollector if one doesn't exist, and adds
// a datapoint to it.
func (self *MonitorGroup) Data(name string, val ...float64) {
	name = SanitizeName(name)
	monitor, err := self.collectors.Get(name, func(_ interface{}) (interface{},
		error) {
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

// Event creates an EventMonitor by the given name if one doesn't exist and
// adds an event to it.
func (self *MonitorGroup) Event(name string) {
	name = SanitizeName(name)
	monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{},
		error) {
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

// Val creates a ValueMonitor by the given name if one doesn't exist and adds
// a value to it.
func (self *MonitorGroup) Val(name string, val float64) {
	name = SanitizeName(name)
	monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{},
		error) {
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

// IntVal is faster than Val when you don't want to deal with floating point
// ops.
func (self *MonitorGroup) IntVal(name string, val int64) {
	name = SanitizeName(name)
	monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{},
		error) {
		return NewIntValueMonitor(), nil
	})
	if err != nil {
		handleError(err)
		return
	}
	val_monitor, ok := monitor.(*IntValueMonitor)
	if !ok {
		handleError(errors.ProgrammerError.New(
			"monitor already exists with different type for name %s", name))
		return
	}
	val_monitor.Add(val)
}

// Chain creates a ChainedMonitor by the given name if one doesn't exist and
// sets the Monitor other to it.
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
