// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"github.com/SpaceMonkeyGo/errors"
	"github.com/SpaceMonkeyGo/monitor/utils"
)

// MonitorStore is a collection of package-level MonitorGroups. There is
// typically only one MonitorStore per process, the DefaultStore.
type MonitorStore struct {
	groups *utils.ThreadsafeCache
}

// NewMonitorStore creates a new MonitorStore
func NewMonitorStore() *MonitorStore {
	return &MonitorStore{groups: utils.NewThreadsafeCache()}
}

// Stats conforms to the Monitor interface
func (s *MonitorStore) Stats(cb func(name string, val float64)) {
	snapshot := s.groups.Snapshot()
	for _, name := range sortedStringKeys(snapshot) {
		cache_val := snapshot[name]
		mon, ok := cache_val.(Monitor)
		if !ok {
			continue
		}
		mon.Stats(cb)
	}
}

// Datapoints conforms to the DataCollection interface
func (s *MonitorStore) Datapoints(reset bool, cb func(name string,
	data [][]float64, total uint64, clipped bool, fraction float64)) {
	snapshot := s.groups.Snapshot()
	for _, name := range sortedStringKeys(snapshot) {
		cache_val := snapshot[name]
		collector, ok := cache_val.(DataCollection)
		if !ok {
			continue
		}
		collector.Datapoints(reset, cb)
	}
}

// GetMonitorsNamed finds or creates a MonitorGroup by the given group name
func (s *MonitorStore) GetMonitorsNamed(group_name string) *MonitorGroup {
	group_name = SanitizeName(group_name)
	cached, err := s.groups.Get(group_name, func(_ interface{}) (interface{}, error) {
		return NewMonitorGroup(group_name), nil
	})
	if err != nil {
		// GetMonitor is often used to initialize global variables, so i'm
		// making an exception for panic
		panic(err)
	}
	group, ok := cached.(*MonitorGroup)
	if !ok {
		// Same
		panic(errors.ProgrammerError.New(
			"non-monitor-group type in monitor group cache!"))
	}
	return group

}

// GetMonitorsNamed finds or creates a MonitorGroup using automatic name
// selection
func (s *MonitorStore) GetMonitors() *MonitorGroup {
	return s.GetMonitorsNamed(PackageName())
}
