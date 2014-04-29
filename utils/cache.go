// Copyright (C) 2013 Space Monkey, Inc.

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
