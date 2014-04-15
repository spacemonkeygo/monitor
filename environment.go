// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"fmt"
	"log"
	"runtime"
	"syscall"

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

	group.Chain("goroutines", MonitorFunc(
		func(cb func(name string, val float64)) {
			cb("count", float64(runtime.NumGoroutine()))
		}))

	group.Chain("memory", MonitorFunc(
		func(cb func(name string, val float64)) {
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)
			MonitorStruct(stats, cb)
		}))

	group.Chain("process", MonitorFunc(
		func(cb func(name string, val float64)) {
			cb("uptime", (space_time.Monotonic() - startTime).Seconds())
			cb("control", 1)
		}))

	group.Chain("runtime", MonitorFunc(
		func(cb func(name string, val float64)) {
			MonitorStruct(RuntimeInternals(), cb)
		}))

	group.Chain("rusage", MonitorFunc(
		func(cb func(name string, val float64)) {
			var rusage syscall.Rusage
			err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)
			if err != nil {
				log.Printf("failed getting rusage data: %s", err)
				return
			}
			MonitorStruct(&rusage, cb)
		}))
}

func SchedulerTrace(out []byte, detailed bool) (n int) {
	var rv int32
	schedTrace(out, detailed, &rv)
	return int(rv)
}

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
