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
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type durationSort []time.Duration

func (d durationSort) Len() int           { return len(d) }
func (d durationSort) Less(i, j int) bool { return d[i] < d[j] }
func (d durationSort) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

// ServeHTTP dumps all of the MonitorStore's keys and values to the requester.
// This method allows a MonitorStore to be registered as an HTTP handler.
func (s *MonitorStore) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if strings.HasSuffix(req.URL.Path, "running") {
		s.Running(func(name string, current []*TaskCtx) {
			fmt.Fprintf(w, "%s - %d tasks\n", name, len(current))
			durs := make([]time.Duration, 0, len(current))
			for _, task := range current {
				durs = append(durs, task.ElapsedTime())
			}
			sort.Sort(sort.Reverse(durationSort(durs)))
			for _, dur := range durs {
				fmt.Fprintf(w, "\t%s\n", dur)
			}
		})
		return
	}

	s.Stats(func(name string, val float64) {
		fmt.Fprintf(w, "%s\t%f\n", name, val)
	})
}
