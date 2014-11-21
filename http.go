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

type sortableTask struct {
	name    string
	elapsed time.Duration
}

type sortableTasks []sortableTask

func (s sortableTasks) Len() int      { return len(s) }
func (s sortableTasks) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s sortableTasks) Less(i, j int) bool {
	return s[i].elapsed < s[j].elapsed
}

// ServeHTTP dumps all of the MonitorStore's keys and values to the requester.
// This method allows a MonitorStore to be registered as an HTTP handler.
func (s *MonitorStore) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if strings.HasSuffix(req.URL.Path, "running") {
		var tasks []sortableTask
		s.Running(func(name string, current []*TaskCtx) {
			for _, task := range current {
				tasks = append(tasks, sortableTask{
					name:    name,
					elapsed: task.ElapsedTime()})
			}
		})
		sort.Sort(sort.Reverse(sortableTasks(tasks)))
		for _, task := range tasks {
			fmt.Fprintf(w, "%s\t%s\n", task.elapsed, task.name)
		}
		return
	}

	if strings.HasSuffix(req.URL.Path, "datapoints") {
		s.Datapoints(false, func(key string, data [][]float64,
			total uint64, clipped bool, fraction float64) {

			fmt.Fprintf(w, "%s\t%d\t%v\t%f\n", key, total, clipped, fraction)
			for idx, points := range data {
				fmt.Fprintf(w, "\t%v", points)
				if (idx+1)%6 == 0 {
					fmt.Fprintln(w)
				}
			}
			fmt.Fprintln(w)
		})
		return
	}

	s.Stats(func(name string, val float64) {
		fmt.Fprintf(w, "%s\t%f\n", name, val)
	})
}
