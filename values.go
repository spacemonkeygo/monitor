// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"math"
	"sync"
)

type ValueMonitor struct {
	mtx         sync.Mutex
	recent      float64
	count       uint64
	sum         float64
	sum_squared float64
	max         float64
	min         float64
}

func NewValueMonitor() *ValueMonitor {
	return &ValueMonitor{
		max: math.Inf(-1),
		min: math.Inf(1)}
}

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
