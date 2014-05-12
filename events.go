// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
