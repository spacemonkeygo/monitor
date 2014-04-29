// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"github.com/SpaceMonkeyGo/monitor/utils"
)

type MonitorGroup struct {
	group_name string
	monitors   *utils.ThreadsafeCache
	collectors *utils.ThreadsafeCache
}

func NewMonitorGroup(name string) *MonitorGroup {
	return &MonitorGroup{
		group_name: SanitizeName(name),
		monitors:   utils.NewThreadsafeCache(),
		collectors: utils.NewThreadsafeCache(),
	}
}
