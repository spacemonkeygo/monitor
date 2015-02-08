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

package monitor

import (
	"sort"

	"github.com/spacemonkeygo/spacelog"
	"gopkg.in/spacemonkeygo/monitor.v1/trace"
)

var (
	// Package-level functions typically work on the DefaultStore. DefaultStore
	// functions as an HTTP handler, if serving statistics over HTTP sounds
	// interesting to you.
	DefaultStore = NewMonitorStore()

	logger = spacelog.GetLogger()

	CallerName             = trace.CallerName
	PackageName            = trace.PackageName
	AddIgnoredCallerPrefix = trace.AddIgnoredCallerPrefix
)

// Monitor is the basic key/value interface. Anything that implements the
// Monitor interface can be connected to the monitor system for later
// processing.
type Monitor interface {
	Stats(cb func(name string, val float64))
}

// RunningTasksCollector keeps track of tasks that are currently in process.
type RunningTasksCollector interface {
	Running(cb func(name string, current []*TaskCtx))
}

// DataCollection is the basic key/vector interface. Anything that implements
// the DataCollection interface can be connected to the monitor system for
// later processing.
type DataCollection interface {
	// Datapoints calls cb with stored datasets. If reset is true, Datapoints
	// should clear its stores. cb is called with the name of the dataset,
	// len(data) datapoints, where a datapoint is a vector of scalars, the total
	// number of datapoints actually seen (which will be >= len(data)), whether
	// or not some datapoints got clipped and the data collector had to revert to
	// stream random sampling, and the fraction of data points collectoed.
	Datapoints(reset bool, cb func(name string, data [][]float64, total uint64,
		clipped bool, fraction float64))
}

func handleError(err error) {
	logger.Errorf("monitoring error: %s", err)
}

func sortedStringKeys(snapshot map[interface{}]interface{}) []string {
	keys := make([]string, 0, len(snapshot))
	for cache_key := range snapshot {
		name, ok := cache_key.(string)
		if !ok {
			continue
		}
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return keys
}

// Stats calls cb with all the statistics registered on the default store.
func Stats(cb func(name string, val float64)) { DefaultStore.Stats(cb) }

// Running calls cb with lists of currently running tasks by name.
func Running(cb func(name string, current []*TaskCtx)) {
	DefaultStore.Running(cb)
}

// Datapoints calls cb with all the datasets registered on the default store.
func Datapoints(reset bool, cb func(name string, data [][]float64, total uint64,
	clipped bool, fraction float64)) {
	DefaultStore.Datapoints(reset, cb)
}

// GetMonitors creates a MonitorGroup with an automatic per-package name on
// the default store.
func GetMonitors() *MonitorGroup {
	return DefaultStore.GetMonitorsNamed(PackageName())
}

// GetMonitorsNamed creates a named MonitorGroup on the default store.
func GetMonitorsNamed(name string) *MonitorGroup {
	return DefaultStore.GetMonitorsNamed(name)
}

// RegisterEnvironment registers environment statistics on the default store.
func RegisterEnvironment() {
	DefaultStore.RegisterEnvironment()
}
