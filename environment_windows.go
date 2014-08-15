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
	"unsafe"
)

var (
	kernel32          = syscall.MustLoadDLL("kernel32.dll")
	getProcCount      = kernel32.MustFindProc("GetProcessHandleCount")
	getModuleFileName = kernel32.MustFindProc("GetModuleFileNameW")
)

func registerPlatformEnvironment(*MonitorGroup) {
}

func fdCount() (count int, err error) {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	var proccount uint32
	r1, _, err := getProcCount.Call(
		uintptr(h), uintptr(unsafe.Pointer(&proccount)))
	if r1 == 0 {
		return 0, err
	}
	return int(proccount), nil
}

func openProc() (*os.File, error) {
	m16 := make([]uint16, 1024)
	r1, _, err := getModuleFileName.Call(0,
		uintptr(unsafe.Pointer(&m16[0])), uintptr(len(m16)))
	if r1 == 0 {
		panic(err)
	}
	m := syscall.UTF16ToString(m16)
	return os.Open(m)
}
