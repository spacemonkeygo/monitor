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
	"math"
	"sync"
)

// ValueMonitor keeps track of the highs and lows and averages and most recent
// versions of some value
type ValueMonitor struct {
	mtx         sync.Mutex
	recent      float64
	count       uint64
	sum         float64
	sum_squared float64
	max         float64
	min         float64
}

// NewValueMonitor creates a new ValueMonitor. You probably want to create a
// new ValueMonitor through MonitorGroup.Val instead.
func NewValueMonitor() *ValueMonitor {
	return &ValueMonitor{
		max: math.Inf(-1),
		min: math.Inf(1)}
}

// Add adds a value to the ValueMonitor
func (v *ValueMonitor) Add(val float64) {
	v.mtx.Lock()
	v.count += 1
	v.sum += val
	v.sum_squared += (val * val)
	v.recent = val
	if val > v.max {
		v.max = val
	}
	if val < v.min {
		v.min = val
	}
	v.mtx.Unlock()
}

// Stats conforms to the Monitor interface
func (v *ValueMonitor) Stats(cb func(name string, val float64)) {
	v.mtx.Lock()
	count := v.count
	sum := v.sum
	sum_squared := v.sum_squared
	recent := v.recent
	max := v.max
	min := v.min
	v.mtx.Unlock()

	if count > 0 {
		cb("avg", sum/float64(count))
	}
	cb("count", float64(count))
	cb("max", max)
	cb("min", min)
	cb("recent", recent)
	cb("sum", sum)
	cb("sum_squared", sum_squared)
}

// IntValueMonitor is faster than ValueMonitor when you don't want to deal with
// floating-point ops
type IntValueMonitor struct {
	mtx         sync.Mutex
	recent      int64
	count       int64
	sum         int64
	sum_squared int64
	max         int64
	min         int64
}

// NewIntValueMonitor returns a new IntValueMonitor. You probably want to
// create a new IntValueMonitor through MonitorGroup.IntVal instead.
func NewIntValueMonitor() *IntValueMonitor {
	return &IntValueMonitor{
		max: math.MinInt64,
		min: math.MaxInt64}
}

// Add adds a value to the IntValueMonitor
func (v *IntValueMonitor) Add(val int64) {
	v.mtx.Lock()
	v.count += 1
	v.sum += val
	v.sum_squared += (val * val)
	v.recent = val
	if val > v.max {
		v.max = val
	}
	if val < v.min {
		v.min = val
	}
	v.mtx.Unlock()
}

// Stats conforms to the Monitor interface
func (v *IntValueMonitor) Stats(cb func(name string, val float64)) {
	v.mtx.Lock()
	count := v.count
	sum := v.sum
	sum_squared := v.sum_squared
	recent := v.recent
	max := v.max
	min := v.min
	v.mtx.Unlock()

	if count > 0 {
		cb("avg", float64(sum/count))
	}
	cb("count", float64(count))
	cb("max", float64(max))
	cb("min", float64(min))
	cb("recent", float64(recent))
	cb("sum", float64(sum))
	cb("sum_squared", float64(sum_squared))
}
