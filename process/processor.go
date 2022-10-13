// author: wsfuyibing <websearch@163.com>
// date: 2022-10-13

package process

import (
	"context"
	"sync"

	"github.com/fuyibing/util/v2/caller"
)

type (
	// Processor
	// 模拟进程接口.
	//
	// 模拟进程生命周期从创建开始到执行结束, 一经结束并销毁.
	Processor interface {
		Healthy() bool

		Restart()
		Start(ctx context.Context) error
		Stop()
	}

	// 模拟进程结构体.
	processor struct {
		mu   sync.RWMutex
		name string

		ci []caller.IgnoreCaller
	}
)

// New
// 创建模拟进程.
func New(name string) Processor {
	return (&processor{name: name}).init()
}

// After
// 注册后置回调.
func (o *processor) After(cs ...caller.IgnoreCaller) Processor {
	o.ci = cs
	return o
}

// Before
// 注册前置回调.
func (o *processor) Before(cs ...caller.IgnoreCaller) Processor {
	o.ci = cs
	return o
}

func (o *processor) Healthy() bool { return false }

func (o *processor) Start(ctx context.Context) error { return nil }
func (o *processor) Stop()                           {}
func (o *processor) Restart()                        {}

// /////////////////////////////////////////////////////////////
// Callback callers execution.
// /////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////
// Instance initializer.
// /////////////////////////////////////////////////////////////

func (o *processor) init() *processor {
	o.mu = sync.RWMutex{}
	return o
}
