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

package utils

import (
	"sync"
)

type ThreadsafeCache struct {
	mtx    sync.RWMutex
	values map[interface{}]interface{}
}

func NewThreadsafeCache() *ThreadsafeCache {
	return &ThreadsafeCache{
		values: make(map[interface{}]interface{}),
	}
}

func (c *ThreadsafeCache) Get(key interface{},
	defaultcb func(key interface{}) (interface{}, error)) (interface{}, error) {

	c.mtx.RLock()
	val, ok := c.values[key]
	if ok {
		c.mtx.RUnlock()
		return val, nil
	}
	c.mtx.RUnlock()

	c.mtx.Lock()
	defer c.mtx.Unlock()

	val, ok = c.values[key]
	if ok {
		return val, nil
	}

	val, err := defaultcb(key)
	if err == nil {
		c.values[key] = val
	}

	return val, err
}

func (c *ThreadsafeCache) Drop(key interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.values, key)
}

func (c *ThreadsafeCache) Snapshot() map[interface{}]interface{} {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	snapshot := make(map[interface{}]interface{}, len(c.values))
	for key, value := range c.values {
		snapshot[key] = value
	}
	return snapshot
}
