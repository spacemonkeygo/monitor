// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
	"flag"
	"sync"

	"code.spacemonkey.com/go/space/rand"
)

var (
	collectionFraction = flag.Float64("monitor.datapoint_collection_fraction", .1,
		"The fraction of datapoints to collect")
	collectionMax = flag.Int("monitor.datapoint_collection_max", 500,
		"The max datapoints to collect")
)

type DatapointCollector struct {
	mtx                 sync.Mutex
	collection_fraction float64
	collection_max      int
	total               uint64
	considered_total    int
	clipped             bool
	dataset             [][]float64
}

func NewDatapointCollector(collection_fraction float64, collection_max int) *DatapointCollector {
	return &DatapointCollector{
		collection_fraction: collection_fraction,
		collection_max:      collection_max}
}

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
