// Copyright (C) 2014 Space Monkey, Inc.

package monitor

import (
	"testing"
)

func ExampleMonitorGroup_Task(t *testing.T) {
	mon := GetMonitors()
	myfunc := func() (err error) {
		defer mon.Task()(&err)
		// do some work
		// maybe return an error
		return nil
	}
	myfunc()
}
