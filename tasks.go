// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"flag"
	"sync"
)

var (
	maxErrorLength = flag.Int("monitor.max_error_length", 40,
		"the max length for an error name")
)

// TaskMonitor is a type for keeping track of tasks. A TaskMonitor will keep
// track of the current number of tasks, the highwater number (the maximum
// amount of concurrent tasks), the total started, the total completed, the
// total that returned without error, the average/min/max/most recent amount
// of time the task took to succeed/fail/both, the number of different kinds
// of errors the task had, and how many times the task had a panic.
//
// N.B.: Error types are best tracked when you're using Space Monkey's
// hierarchical error package: http://github.com/SpaceMonkeyGo/errors
type TaskMonitor struct {
	mtx             sync.Mutex
	current         uint64
	highwater       uint64
	total_started   uint64
	total_completed uint64
	success         uint64
	success_timing  *IntValueMonitor
	error_timing    *IntValueMonitor
	total_timing    *IntValueMonitor
	errors          map[string]uint64
	panics          uint64
}

// NewTaskMonitor returns a new TaskMonitor. You probably want to create
// a TaskMonitor using MonitorGroup.Task instead.
func NewTaskMonitor() *TaskMonitor {
	return &TaskMonitor{
		errors:         make(map[string]uint64),
		success_timing: NewIntValueMonitor(),
		error_timing:   NewIntValueMonitor(),
		total_timing:   NewIntValueMonitor()}
}
