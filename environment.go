// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
	"runtime"

	space_time "code.spacemonkey.com/go/space/time"
)

var (
	startTime = space_time.Monotonic()
)

func (store *MonitorStore) RegisterEnvironment() {
	if store == nil {
		store = DefaultStore
	}
	group := store.GetMonitorsNamed("env")
	group.Chain("goroutines", GoroutineMonitor{})
	group.Chain("memory", MemoryMonitor{})
	group.Chain("process", ProcessMonitor{})
}

type GoroutineMonitor struct{}

func (GoroutineMonitor) Stats(cb func(name string, val float64)) {
	cb("count", float64(runtime.NumGoroutine()))
}

type MemoryMonitor struct{}

func (MemoryMonitor) Stats(cb func(name string, val float64)) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	MonitorStruct(stats, cb)
}

type ProcessMonitor struct{}

func (ProcessMonitor) Stats(cb func(name string, val float64)) {
	cb("uptime", (space_time.Monotonic() - startTime).Seconds())
	cb("control", 1)
}
