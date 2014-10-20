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

package trace_test

import (
	"github.com/spacemonkeygo/monitor/trace"
)

func ExampleSpan_Observe() {
	myfunc := func() (err error) {
		span := trace.NewTrace("my span")
		defer span.Observe()(&err)

		// do things

		return nil
	}

	myfunc()
}
