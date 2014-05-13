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

	"github.com/spacemonkeygo/crc"
)

// MonitorFunc assists in Monitor interface instances
type MonitorFunc func(cb func(name string, val float64))

// Stats just calls f with the given cb
func (f MonitorFunc) Stats(cb func(name string, val float64)) { f(cb) }

// PrefixStats will call cb with all of the same calls obj would have, except
// every name is prefixed with name.
func PrefixStats(name string, obj Monitor, cb func(name string, val float64)) {
	obj.Stats(func(sub_name string, val float64) {
		cb(fmt.Sprintf("%s.%s", name, sub_name), val)
	})
}

// Collect takes something that implements the Monitor interface and returns
// a key/value map.
func Collect(mon Monitor) map[string]float64 {
	rv := make(map[string]float64)
	mon.Stats(func(name string, val float64) {
		rv[name] = val
	})
	return rv
}

// BoolAsFloat converts a bool value into a float64 value, for easier datapoint
// usage.
func BoolAsFloat(val bool) float64 {
	if val {
		return 1
	}
	return 0
}

// FloatHash CRCs a byte array and converts the CRC into a float64, for easier
// datapoint usage.
func FloatHash(data []byte) float64 {
	return float64(crc.CRC(crc.InitialCRC, data))
}

// SanitizeName cleans a stat or datapoint name to be representable in a wide
// range of data collection software.
func SanitizeName(name string) string {
	rname := []byte(name)
	for i, r := range rname {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		default:
			switch r {
			case '_', '.', '-':
			case '/':
				rname[i] = '.'
			default:
				rname[i] = '_'
			}
		}
	}
	return string(rname)
}
