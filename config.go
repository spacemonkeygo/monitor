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

// Config is a configuration struct meant to be used with
//   github.com/spacemonkeygo/flagfile/utils.Setup
// but can be set independently.
var Config = struct {
	DefaultCollectionFraction float64 `default:".1" usage:"The fraction of datapoints to collect"`
	DefaultCollectionMax      int     `default:"500" usage:"The max datapoints to collect"`
	MaxErrorLength            int     `default:"40" usage:"the max length for an error name"`
}{
	DefaultCollectionFraction: .1,
	DefaultCollectionMax:      500,
	MaxErrorLength:            40,
}
