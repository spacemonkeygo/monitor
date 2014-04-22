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

func NewTaskMonitor() *TaskMonitor {
	return &TaskMonitor{
		errors:         make(map[string]uint64),
		success_timing: NewIntValueMonitor(),
		error_timing:   NewIntValueMonitor(),
		total_timing:   NewIntValueMonitor()}
}
