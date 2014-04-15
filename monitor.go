// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"runtime"
	"sort"
	"strings"

	space_log "code.spacemonkey.com/go/space/log"
)

var (
	DefaultStore = NewMonitorStore()

	logger = space_log.GetLogger()
)

type Monitor interface {
	Stats(cb func(name string, val float64))
}

type DataCollection interface {
	Datapoints(reset bool, cb func(name string, data [][]float64, total uint64,
		clipped bool, fraction float64))
}

func CallerName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown.unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown.unknown"
	}
	return strings.TrimPrefix(f.Name(), "code.spacemonkey.com/go/")
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

func Stats(cb func(name string, val float64)) { DefaultStore.Stats(cb) }

func Datapoints(reset bool, cb func(name string, data [][]float64, total uint64,
	clipped bool, fraction float64)) {
	DefaultStore.Datapoints(reset, cb)
}

func GetMonitors() *MonitorGroup {
	return DefaultStore.GetMonitorsNamed(CallerName())
}

func GetMonitorsNamed(name string) *MonitorGroup {
	return DefaultStore.GetMonitorsNamed(name)
}

func RegisterEnvironment() {
	DefaultStore.RegisterEnvironment()
}
