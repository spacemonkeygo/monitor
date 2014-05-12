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

// +build no_mon

package monitor

func (g *MonitorGroup) Stats(cb func(name string, val float64)) {}

func (g *MonitorGroup) Datapoints(reset bool, cb func(name string,
	data [][]float64, total uint64, clipped bool, fraction float64)) {
}

func (self *MonitorGroup) Data(name string, val ...float64) {}
func (self *MonitorGroup) Event(name string)                {}
func (self *MonitorGroup) Val(name string, val float64)     {}
func (self *MonitorGroup) IntVal(name string, val int64)    {}
func (self *MonitorGroup) Chain(name string, other Monitor) {}

func (self *MonitorGroup) Task() func(*error)     { return func(*error) {} }
func (self *MonitorGroup) DataTask() func(*error) { return func(*error) {} }

func (self *MonitorGroup) TaskNamed(name string) func(*error) {
	return func(*error) {}
}
