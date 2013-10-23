// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "flag"
    "math/rand"
    "sync"

    "code.spacemonkey.com/go/errors"
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

func (self *MonitorGroup) Data(name string, val ...float64) {
    monitor, err := self.monitors.Get(name, func(_ interface{}) (interface{}, error) {
        return NewDatapointCollector(*collectionFraction, *collectionMax), nil
    })
    if err != nil {
        handleError(err)
        return
    }
    datapoint_collector, ok := monitor.(*DatapointCollector)
    if !ok {
        handleError(errors.ProgrammerError.New(
            "monitor already exists with different type for name %s", name))
        return
    }
    datapoint_collector.Add(val...)
}
