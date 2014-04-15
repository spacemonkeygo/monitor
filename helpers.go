// Copyright (C) 2013 Space Monkey, Inc.

package monitor

import (
	"fmt"
	"regexp"

	"code.spacemonkey.com/go/space/crc"
	"code.spacemonkey.com/go/types"
)

var (
	badChars = regexp.MustCompile("[^a-zA-Z0-9_.-]")
	slashes  = regexp.MustCompile("[/]")
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
	return badChars.ReplaceAllLiteralString(
		slashes.ReplaceAllLiteralString(name, "."), "_")
}
