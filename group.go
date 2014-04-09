// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
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
