// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"github.com/SpaceMonkeyGo/monitor/utils"
)

// MonitorGroup is a collection of named Monitor interfaces and DataCollector
// interfaces. They are automatically created by various calls on the
// MonitorGroup
type MonitorGroup struct {
	group_name string
	monitors   *utils.ThreadsafeCache
	collectors *utils.ThreadsafeCache
}

// NewMonitorGroup makes a new MonitorGroup unattached to anything.
func NewMonitorGroup(name string) *MonitorGroup {
	return &MonitorGroup{
		group_name: SanitizeName(name),
		monitors:   utils.NewThreadsafeCache(),
		collectors: utils.NewThreadsafeCache(),
	}
}
