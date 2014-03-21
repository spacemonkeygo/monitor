// Copyright (C) 2013 Space Monkey, Inc.

// +build no_mon

package client

func (g *MonitorGroup) Stats(cb func(name string, val float64)) {}

func (g *MonitorGroup) Datapoints(reset bool, cb func(name string,
	data [][]float64, total uint64, clipped bool, fraction float64)) {
}

func (self *MonitorGroup) Data(name string, val ...float64) {}
func (self *MonitorGroup) Event(name string)                {}
func (self *MonitorGroup) Val(name string, val float64)     {}
func (self *MonitorGroup) Chain(name string, other Monitor) {}

func (self *MonitorGroup) Task() func(*error) {
	return func(*error) {}
}

func (self *MonitorGroup) TaskNamed(name string) func(*error) {
	return func(*error) {}
}
