// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"fmt"
	"log"
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
	group.Chain("runtime", RuntimeMonitor{})
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

type RuntimeMonitor struct{}

func (RuntimeMonitor) Stats(cb func(name string, val float64)) {
	MonitorStruct(RuntimeInternals(), cb)
}

func SchedulerTrace(out []byte, detailed bool) (n int) {
	var rv int32
	schedTrace(out, detailed, &rv)
	return int(rv)
}

// InternalStats
// shared with C. If you edit this struct, edit IStats
type InternalStats struct {
	GoMaxProcs  int32
	IdleProcs   int32
	ThreadCount int32
	IdleThreads int32
	RunQueue    int32
}

func schedTrace(b []byte, detailed bool, n *int32)

func schedTraceData(stats *InternalStats) {
	var data [256]byte
	SchedulerTrace(data[:], false)
	var uptime int64
	n, err := fmt.Sscanf(string(data[:]),
		"SCHED %dms: gomaxprocs=%d idleprocs=%d threads=%d idlethreads=%d "+
			"runqueue=%d", &uptime, &stats.GoMaxProcs, &stats.IdleProcs,
		&stats.ThreadCount, &stats.IdleThreads, &stats.RunQueue)
	if err != nil || n != 6 {
		log.Printf("failed getting runtime data from scheduler trace: %v, %d",
			err, n)
	}
}

func RuntimeInternals() (rv InternalStats) {
	schedTraceData(&rv)
	return rv
}
