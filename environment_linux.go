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
	"os"
	"syscall"
)

func registerPlatformEnvironment(group *MonitorGroup) {
	group.Chain("rusage", MonitorFunc(
		func(cb func(name string, val float64)) {
			var rusage syscall.Rusage
			err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)
			if err != nil {
				logger.Errorf("failed getting rusage data: %s", err)
				return
			}
			MonitorStruct(&rusage, cb)
		}))
}

func fdCount() (count int, err error) {
	f, err := os.Open("/proc/self/fd")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	for {
		names, err := f.Readdirnames(4096)
		count += len(names)
		if err != nil {
			if err == io.EOF {
				return count, nil
			}
			return count, err
		}
	}
}

func openProc() (*os.File, error) {
	return os.Open("/proc/self/exe")
}
