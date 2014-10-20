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
	"io"
	"runtime"

	"github.com/spacemonkeygo/crc"
	"github.com/spacemonkeygo/monotime"
)

var (
	startTime = monotime.Monotonic()
)

// RegisterEnvironment configures the MonitorStore receiver to understand all
// sorts of process environment statistics, such as memory statistics,
// process uptime, file descriptor use, goroutine use, runtime internals,
// Rusage stats, etc.
func (store *MonitorStore) RegisterEnvironment() {
	if store == nil {
		store = DefaultStore
	}
	group := store.GetMonitorsNamed("env")

	group.Chain("goroutines", MonitorFunc(
		func(cb func(name string, val float64)) {
			cb("count", float64(runtime.NumGoroutine()))
		}))

	group.Chain("memory", MonitorFunc(
		func(cb func(name string, val float64)) {
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)
			MonitorStruct(stats, cb)
		}))

	process_crc, err := ProcessCRC()
	if err != nil {
		logger.Errorf("failed determining process crc: %s", err)
		process_crc = 0
	}

	group.Chain("process", MonitorFunc(
		func(cb func(name string, val float64)) {
			cb("control", 1)
			fds, err := FdCount()
			if err != nil {
				logger.Errorf("failed getting fd count: %s", err)
			} else {
				cb("fds", float64(fds))
			}
			cb("crc", float64(process_crc))
			cb("uptime", (monotime.Monotonic() - startTime).Seconds())
		}))

	group.Chain("runtime", MonitorFunc(
		func(cb func(name string, val float64)) {
			MonitorStruct(RuntimeInternals(), cb)
		}))

	registerPlatformEnvironment(group)
}

// FdCount counts how many open file descriptors there are.
func FdCount() (count int, err error) {
	return fdCount()
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) { return f(p) }

func ProcessCRC() (uint32, error) {
	fh, err := openProc()
	if err != nil {
		return 0, err
	}
	defer fh.Close()
	c := crc.InitialCRC
	_, err = io.Copy(writerFunc(func(p []byte) (n int, err error) {
		c = crc.CRC(c, p)
		return len(p), nil
	}), fh)
	if err != nil {
		return 0, err
	}
	return c, nil
}
