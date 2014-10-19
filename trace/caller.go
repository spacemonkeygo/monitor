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

package trace

import (
	"runtime"
	"strings"
)

var (
	// IgnoredPrefixes is a list of prefixes to ignore when performing automatic
	// name generation.
	ignoredPrefixes []string
)

// AddIgnoredCallerPrefix adds a prefix that will get trimmed off from
// automatic name generation through CallerName (and therefore Trace)
func AddIgnoredCallerPrefix(prefix string) {
	ignoredPrefixes = append(ignoredPrefixes, prefix)
}

// CallerName returns the name of the caller two frames up the stack.
func CallerName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown.unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown.unknown"
	}
	name := f.Name()
	for _, prefix := range ignoredPrefixes {
		name = strings.TrimPrefix(name, prefix)
	}
	return name
}

// PackageName returns the name of the package of the caller two frames up
// the stack.
func PackageName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown"
	}
	name := f.Name()
	for _, prefix := range ignoredPrefixes {
		name = strings.TrimPrefix(name, prefix)
	}
	return name

	idx := strings.Index(name, ".")
	if idx >= 0 {
		name = name[:idx]
	}
	return name
}
