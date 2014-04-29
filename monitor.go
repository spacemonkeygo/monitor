// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"runtime"
	"sort"
	"strings"

	"github.com/SpaceMonkeyGo/spacelog"
)

var (
	// Package-level functions typically work on the DefaultStore. DefaultStore
	// functions as an HTTP handler, if serving statistics over HTTP sounds
	// interesting to you.
	DefaultStore = NewMonitorStore()

	// IgnoredPrefixes is a list of prefixes to ignore when performing automatic
	// name generation.
	IgnoredPrefixes []string

	logger = spacelog.GetLogger()
)

// Monitor is the basic key/value interface. Anything that implements the
// Monitor interface can be connected to the monitor system for later
// processing.
type Monitor interface {
	Stats(cb func(name string, val float64))
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

// CallerName returns the name of the caller two frames up the stack.
func CallerName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown.unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown.unknown"
	}
	name := f.Name()
	for _, prefix := range IgnoredPrefixes {
		name = strings.TrimPrefix(name, prefix)
	}
	return name
}

// PackageName returns the name of the package of the caller two frames up
// the stack.
func PackageName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown"
	}
	name := f.Name()
	for _, prefix := range IgnoredPrefixes {
		name = strings.TrimPrefix(name, prefix)
	}
	return name

	idx := strings.Index(name, ".")
	if idx >= 0 {
		name = name[:idx]
	}
	return name
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
