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

/*
Package monitor is a flexible code instrumenting and data collection library.

With this package, it's easy to monitor and watch all sorts of data.
A motivating example:

	package main

	import (
		"net/http"

		"github.com/spacemonkeygo/monitor"
	)

	var (
		mon = monitor.GetMonitors()
	)

	func FixSerenity() (err error) {
		defer mon.Task()(&err)

		if SerenityBroken() {
			err := CallKaylee()
			mon.Event("kaylee called")
			if err != nil {
				return err
			}
		}

		stowaway_count := StowawaysNeedingHiding()
		mon.Val("stowaway count", stowaway_count)
		err = HideStowaways(stowaway_count)
		if err != nil {
			return err
		}

		return nil
	}

	func Monitor() {
		go http.ListenAndServe(":8080", monitor.DefaultStore)
	}

In this example, calling FixSerenity will cause the endpoint at
http://localhost:8080/ to return all sorts of data, such as:

 * How many times we've needed to fix the Serenity
   (the Task monitor infers the statistic name from the callstack)
 * How many times we've succeeded
 * How many times we've failed
 * How long it's taken each time (min/max/avg/recent)
 * How many times we needed to call Kaylee
 * How many errors we've received (per error type!)
 * Statistics on how many stowaways we usually have (min/max/avg/recent/etc)

To collect these statistics without the http server, you can call
monitor.Stats like so

	monitor.Stats(func(name string, val float64) {
		// do something with name, val
	})

This package lets you easily instrument your code with all of these goodies and
more!
*/
package monitor
