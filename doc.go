// Copyright (C) 2014 Space Monkey, Inc.

/*
Package monitor is a flexible code instrumenting and data collection library.

With this package, it's easy to monitor and watch all sorts of data.
A motivating example:

	package main

	import (
		"net/http"

		"github.com/SpaceMonkeyGo/monitor"
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

This package lets you easily instrument your code with all of these goodies and
more!
*/
package monitor
