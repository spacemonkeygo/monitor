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

package sched

import (
	"fmt"

	"github.com/spacemonkeygo/spacelog"
)

var (
	logger = spacelog.GetLogger()
)

func schedTrace(b []byte, detailed bool, n *int32)

// SchedulerTrace collects the output of the standard scheduler trace debug
// output line
func SchedulerTrace(out []byte, detailed bool) (n int) {
	var rv int32
	schedTrace(out, detailed, &rv)
	return int(rv)
}

// InternalStats represents the data typically displayed in a scheduler trace.
type InternalStats struct {
	GoMaxProcs  int32
	IdleProcs   int32
	ThreadCount int32
	IdleThreads int32
	RunQueue    int32
}

func schedTraceData(stats *InternalStats) {
	var data [256]byte
	SchedulerTrace(data[:], false)
	var uptime int64
	n, err := fmt.Sscanf(string(data[:]),
		"SCHED %dms: gomaxprocs=%d idleprocs=%d threads=%d idlethreads=%d "+
			"runqueue=%d", &uptime, &stats.GoMaxProcs, &stats.IdleProcs,
		&stats.ThreadCount, &stats.IdleThreads, &stats.RunQueue)
	if err != nil || n != 6 {
		logger.Errorf("failed getting runtime data from scheduler trace: %v, %d",
			err, n)
	}
}

// RuntimeInternals parses a scheduler trace line into an InternalStats struct
func RuntimeInternals() (rv InternalStats) {
	schedTraceData(&rv)
	return rv
}
