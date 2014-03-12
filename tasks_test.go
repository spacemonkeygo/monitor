// Copyright (C) 2014 Space Monkey, Inc.

package client

import (
	"io"
	"strings"
	"testing"
)

func check(t *testing.T, mon Monitor, success, total, errors, panics float64) {
	var found_errors float64
	mon.Stats(func(name string, val float64) {
		switch {
		case name == "success" && val != success:
			t.Errorf("unexpected success count: %f != %f", val, success)
		case name == "total" && val != total:
			t.Errorf("unexpected success count: %f != %f", val, total)
		case name == "panics" && val != panics:
			t.Errorf("unexpected panics count: %f != %f", val, panics)
		case strings.HasPrefix(name, "error_"):
			found_errors += val
		}
	})
	if found_errors != errors {
		t.Errorf("unexpected errors count: %f != %f", found_errors, errors)
	}
}

func ignore([]byte) {}

func TestTaskPanics(t *testing.T) {
	mon := NewTaskMonitor()
	check(t, mon, 0, 0, 0, 0)

	func() {
		defer mon.Start().Finish(nil)
	}()
	check(t, mon, 1, 1, 0, 0)

	func() {
		var err error
		defer mon.Start().Finish(&err)
		err = io.EOF
	}()
	check(t, mon, 1, 2, 1, 0)

	func() {
		defer func() { recover() }()
		func() {
			var list []byte
			defer mon.Start().Finish(nil)
			ignore(list[4:7])
		}()
	}()
	check(t, mon, 1, 3, 2, 1)

	func() {
		defer func() { recover() }()
		func() {
			var nilref *testing.T
			defer mon.Start().Finish(nil)
			nilref.Fatalf("this should fail")
		}()
	}()
	check(t, mon, 1, 4, 3, 2)

	func() {
		defer func() { recover() }()
		func() {
			defer mon.Start().Finish(nil)
			panic("waaah")
		}()
	}()
	check(t, mon, 1, 5, 4, 3)

	func() {
		defer mon.Start().Finish(nil)
	}()
	check(t, mon, 2, 6, 4, 3)
}
