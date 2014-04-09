// Copyright (C) 2013 Space Monkey, Inc.

// +build no_mon

package monitor

func (t *TaskMonitor) Stats(cb func(name string, val float64)) {}

func (t *TaskMonitor) Start() func(*error)  { return func(*error) {} }
func (t *TaskMonitor) NewContext() *TaskCtx { return &TaskCtx{} }

type TaskCtx struct{}

func (c *TaskCtx) Finish(err_ref *error, rec interface{}) {}
