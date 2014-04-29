// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"flag"
	"math/rand"
	"sync"
)

var (
	collectionFraction = flag.Float64("monitor.datapoint_collection_fraction", .1,
		"The fraction of datapoints to collect")
	collectionMax = flag.Int("monitor.datapoint_collection_max", 500,
		"The max datapoints to collect")
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

	if rand.Float64() >= d.collection_fraction {
		return
	}

	d.considered_total += 1
	if d.clipped {
		r := rand.Intn(d.considered_total)
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
