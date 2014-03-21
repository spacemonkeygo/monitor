// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
	"sync"
)

type EventMonitor struct {
	mtx   sync.Mutex
	count uint64
}

func NewEventMonitor() *EventMonitor {
	return &EventMonitor{}
}

func (e *EventMonitor) Add() {
	e.mtx.Lock()
	e.count += 1
	e.mtx.Unlock()
}

func (e *EventMonitor) Stats(cb func(name string, val float64)) {
	e.mtx.Lock()
	count := e.count
	e.mtx.Unlock()
	cb("count", float64(count))
}
