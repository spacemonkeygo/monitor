// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"fmt"
	"net/http"
)

// ServeHTTP dumps all of the MonitorStore's keys and values to the requester.
// This method allows a MonitorStore to be registered as an HTTP handler.
func (s *MonitorStore) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	s.Stats(func(name string, val float64) {
		fmt.Fprintf(w, "%s\t%f\n", name, val)
	})
}
