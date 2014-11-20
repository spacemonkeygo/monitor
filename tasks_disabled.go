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

// +build no_mon

package monitor

func (t *TaskMonitor) Stats(cb func(name string, val float64)) {}

func (t *TaskMonitor) Start() func(*error)  { return func(*error) {} }
func (t *TaskMonitor) NewContext() *TaskCtx { return &TaskCtx{} }

func (c *TaskCtx) Finish(err_ref *error, rec interface{}) {}

func (t *TaskMonitor) Running() (rv []*TaskCtx) { return nil }
