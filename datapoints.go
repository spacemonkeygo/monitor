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
	"sync"

	"gopkg.in/spacemonkeygo/monitor.v1/trace"
)

// DatapointCollector collects a set of datapoints
type DatapointCollector struct {
	mtx                 sync.Mutex
	collection_fraction float64
	collection_max      int
	total               uint64
	considered_total    int
	clipped             bool
	dataset             [][]float64
}

// NewDatapointCollector makes a new DatapointCollector that will collect
// collection_fraction of all datapoints, and will start clipping data once
// collection_max has been reached without getting drained.
//
// You probably want to create a new DatapointCollector using MonitorGroup.Data
// instead.
func NewDatapointCollector(collection_fraction float64, collection_max int) *DatapointCollector {
	return &DatapointCollector{
		collection_fraction: collection_fraction,
		collection_max:      collection_max}
}

// Add adds new datapoints to the collector. A datapoint is an n-dimensional
// vector. There should be one argument for each scalar in the vector.
func (d *DatapointCollector) Add(val ...float64) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	d.total += 1

	if trace.Rng.Float64() >= d.collection_fraction {
		return
	}

	d.considered_total += 1
	if d.clipped {
		r := trace.Rng.Intn(d.considered_total)
		if r < len(d.dataset) {
			d.dataset[r] = val
		}
	} else {
		d.dataset = append(d.dataset, val)
		if len(d.dataset) >= d.collection_max {
			d.clipped = true
		}
	}
}

// Datapoints returns all of the saved datapoints and any statistics about the
// dataset retained to cb. If reset is true, the collector will be reset and
// the datapoints will be drained from the collector. Datapoints conforms to
// the DataCollection interface.
func (d *DatapointCollector) Datapoints(reset bool, cb func(name string,
	data [][]float64, total uint64, clipped bool, fraction float64)) {

	d.mtx.Lock()
	total := d.total
	clipped := d.clipped
	fraction := d.collection_fraction
	var data_out [][]float64
	if reset {
		data_out = d.dataset
		d.dataset = nil
		d.total = 0
		d.clipped = false
		d.considered_total = 0
	} else {
		data_out = make([][]float64, 0, len(d.dataset))
		for _, row := range d.dataset {
			data_out = append(data_out, row)
		}
	}
	d.mtx.Unlock()

	cb("data", data_out, total, clipped, fraction)
}
