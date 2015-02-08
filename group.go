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
	"gopkg.in/spacemonkeygo/monitor.v1/utils"
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
