// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"fmt"

	"code.spacemonkey.com/go/space/crc"
	"code.spacemonkey.com/go/types"
)

type MonitorFunc func(cb func(name string, val float64))

func (f MonitorFunc) Stats(cb func(name string, val float64)) { f(cb) }

func PrefixStats(name string, obj Monitor, cb func(name string, val float64)) {
	obj.Stats(func(sub_name string, val float64) {
		cb(fmt.Sprintf("%s.%s", name, sub_name), val)
	})
}

func Collect(mon Monitor) map[string]float64 {
	rv := make(map[string]float64)
	mon.Stats(func(name string, val float64) {
		rv[name] = val
	})
	return rv
}

func BoolAsFloat(val bool) float64 {
	if val {
		return 1
	}
	return 0
}

func FloatHash(data types.Binary) float64 {
	return float64(crc.CRC(crc.InitialCRC, []byte(data)))
}

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
