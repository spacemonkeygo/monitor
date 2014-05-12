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
	"strings"
	"testing"
)

func check(t *testing.T, mon Monitor, success, total, errors, panics float64) {
	var found_errors float64
	mon.Stats(func(name string, val float64) {
		switch {
		case name == "foo.bar.success" && val != success:
			t.Errorf("unexpected success count: %f != %f", val, success)
		case name == "foo.bar.total" && val != total:
			t.Errorf("unexpected success count: %f != %f", val, total)
		case name == "foo.bar.panics" && val != panics:
			t.Errorf("unexpected panics count: %f != %f", val, panics)
		case strings.HasPrefix(name, "foo.bar.error_"):
			found_errors += val
		default:
		}
	})
	if found_errors != errors {
		t.Errorf("unexpected errors count: %f != %f", found_errors, errors)
	}
}

func ignore([]byte) {}

func TestTaskPanics(t *testing.T) {
	mon := NewMonitorGroup("foo")
	check(t, mon, 0, 0, 0, 0)

	func() {
		defer mon.TaskNamed("bar")(nil)
	}()
	check(t, mon, 1, 1, 0, 0)

	func() {
		var err error
		defer mon.TaskNamed("bar")(&err)
		err = io.EOF
	}()
	check(t, mon, 1, 2, 1, 0)

	func() {
		defer func() { recover() }()
		func() {
			var list []byte
			defer mon.TaskNamed("bar")(nil)
			ignore(list[4:7])
		}()
	}()
	check(t, mon, 1, 3, 2, 1)

	func() {
		defer func() { recover() }()
		func() {
			var nilref *testing.T
			defer mon.TaskNamed("bar")(nil)
			nilref.Fatalf("this should fail")
		}()
	}()
	check(t, mon, 1, 4, 3, 2)

	func() {
		defer func() { recover() }()
		func() {
			defer mon.TaskNamed("bar")(nil)
			panic("waaah")
		}()
	}()
	check(t, mon, 1, 5, 4, 3)

	func() {
		defer mon.TaskNamed("bar")(nil)
	}()
	check(t, mon, 2, 6, 4, 3)
}

func ExampleTaskMonitor_Start(t *testing.T) {
	task := NewTaskMonitor()
	myfunc := func() (err error) {
		defer task.Start()(&err)
		// do some work
		// maybe return an error
		return nil
	}
	myfunc()
}
