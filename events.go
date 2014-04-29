// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"sync"
)

// EventMonitor keeps track of the number of times an event happened
type EventMonitor struct {
	mtx   sync.Mutex
	count uint64
}

// NewEventMonitor makes a new event monitor. You probably want to create a
// new EventMonitor using MonitorGroup.Event instead.
func NewEventMonitor() *EventMonitor {
	return &EventMonitor{}
}

// Add indicates that the given event happened again
func (e *EventMonitor) Add() {
	e.mtx.Lock()
	e.count += 1
	e.mtx.Unlock()
}

// Stats conforms to the Monitor interface
func (e *EventMonitor) Stats(cb func(name string, val float64)) {
	e.mtx.Lock()
	count := e.count
	e.mtx.Unlock()
	cb("count", float64(count))
}
