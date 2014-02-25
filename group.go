// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
	"fmt"

	"code.spacemonkey.com/go/space/sync"
)

type MonitorGroup struct {
	group_name string
	monitors   *sync.ThreadsafeCache
	collectors *sync.ThreadsafeCache
}

func NewMonitorGroup(name string) *MonitorGroup {
	return &MonitorGroup{
		group_name: SanitizeName(name),
		monitors:   sync.NewThreadsafeCache(),
		collectors: sync.NewThreadsafeCache(),
	}
}

func (g *MonitorGroup) Stats(cb func(name string, val float64)) {
	snapshot := g.monitors.Snapshot()
	for _, name := range sortedStringKeys(snapshot) {
		cache_val := snapshot[name]
		mon, ok := cache_val.(Monitor)
		if !ok {
			continue
		}
		mon.Stats(func(subname string, val float64) {
			cb(fmt.Sprintf("%s.%s.%s", g.group_name, name, subname), val)
		})
	}
}

func (g *MonitorGroup) Datapoints(reset bool, cb func(name string,
	data [][]float64, total uint64, clipped bool, fraction float64)) {
	snapshot := g.collectors.Snapshot()
	for _, name := range sortedStringKeys(snapshot) {
		cache_val := snapshot[name]
		collector, ok := cache_val.(DataCollection)
		if !ok {
			continue
		}
		collector.Datapoints(reset, func(subname string, data [][]float64,
			total uint64, clipped bool, fraction float64) {
			cb(fmt.Sprintf("%s.%s.%s", g.group_name, name, subname), data,
				total, clipped, fraction)
		})
	}
}
