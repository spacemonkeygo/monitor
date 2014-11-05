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
	"testing"
)

func TestReadProcStatSelf(t *testing.T) {
	var stat procSelfStat
	err := readProcSelfStat(&stat)
	if err != nil {
		t.Fatal(err)
	}
	data := make(map[string]float64)
	MonitorStruct(&stat, func(key string, val float64) { data[key] = val })

	for _, key := range []string{"Blocked", "Endcode", "Nice", "Vsize", "Rss"} {
		_, exists := data[key]
		if !exists {
			t.Fatalf("%s doesn't exist", key)
		}
	}
}
