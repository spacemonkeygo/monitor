// Copyright (C) 2013 Space Monkey, Inc.

package client

import (
    "fmt"
    "net/http"
)

func (s *MonitorStore) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    s.Stats(func(name string, val float64) {
        fmt.Fprintf(w, "%s\t%f\n", name, val)
    })
}
