// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"

	"github.com/SpaceMonkeyGo/crc"
	"github.com/SpaceMonkeyGo/monotime"
)

var (
	startTime = monotime.Monotonic()
)

// RegisterEnvironment configures the MonitorStore receiver to understand all
// sorts of process environment statistics, such as memory statistics,
// process uptime, file descriptor use, goroutine use, runtime internals,
// Rusage stats, etc.
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

	process_crc, err := ProcessCRC()
	if err != nil {
		logger.Errorf("failed determining process crc: %s", err)
		process_crc = 0
	}

	group.Chain("process", MonitorFunc(
		func(cb func(name string, val float64)) {
			cb("control", 1)
			fds, err := FdCount()
			if err != nil {
				logger.Errorf("failed getting fd count: %s", err)
			} else {
				cb("fds", float64(fds))
			}
			cb("crc", float64(process_crc))
			cb("uptime", (monotime.Monotonic() - startTime).Seconds())
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
				logger.Errorf("failed getting rusage data: %s", err)
				return
			}
			MonitorStruct(&rusage, cb)
		}))
}

// SchedulerTrace collects the output of the standard scheduler trace debug
// output line
func SchedulerTrace(out []byte, detailed bool) (n int) {
	var rv int32
	schedTrace(out, detailed, &rv)
	return int(rv)
}

// InternalStats represents the data typically displayed in a scheduler trace.
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
		logger.Errorf("failed getting runtime data from scheduler trace: %v, %d",
			err, n)
	}
}

// RuntimeInternals parses a scheduler trace line into an InternalStats struct
func RuntimeInternals() (rv InternalStats) {
	schedTraceData(&rv)
	return rv
}

// FdCount counts how many open file descriptors there are.
func FdCount() (count int, err error) {
	f, err := os.Open("/proc/self/fd")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	for {
		names, err := f.Readdirnames(4096)
		count += len(names)
		if err != nil {
			if err == io.EOF {
				return count, nil
			}
			return count, err
		}
	}
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) { return f(p) }

func ProcessCRC() (uint32, error) {
	fh, err := os.Open("/proc/self/exe")
	if err != nil {
		return 0, err
	}
	defer fh.Close()
	c := crc.InitialCRC
	_, err = io.Copy(writerFunc(func(p []byte) (n int, err error) {
		c = crc.CRC(c, p)
		return len(p), nil
	}), fh)
	if err != nil {
		return 0, err
	}
	return c, nil
}
