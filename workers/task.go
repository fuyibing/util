// author: wsfuyibing <websearch@163.com>
// date: 2022-11-16

package workers

import (
	"context"
	"fmt"
	"github.com/fuyibing/log/v3"
	"sync"
	"sync/atomic"
	"time"
)

var (
	taskId   uint64
	taskPool sync.Pool
)

type (
	// Task
	// 任务接口.
	Task interface {
		Run() (success bool)
		SetContext(ctx context.Context) Task
		SetFinish(finish TaskFinish) Task
		SetHandler(handler TaskHandler) Task
	}

	// TaskFinish
	// 任务完成回调.
	TaskFinish func(ctx context.Context, res TaskResult)

	// TaskHandler
	// 任务处理回调.
	TaskHandler func(ctx context.Context) interface{}

	// 任务结构体.
	task struct {
		acquired     time.Time
		acquires, id uint64

		ctx     context.Context
		finish  TaskFinish
		handler TaskHandler
	}
)

// NewTask
// 从池中取出实例.
func NewTask() Task {
	return taskPool.Get().(*task).before()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *task) Run() (success bool) {
	// 准备执行.
	res := NewTaskResult()
	res.Begin = time.Now()

	// 执行任务.
	if o.handler == nil {
		res.Finish = time.Now()
		res.Errors = append(res.Errors, fmt.Errorf("handler callback not defined"))
	} else {
		success = func(ctx context.Context, callback TaskHandler) (suc bool) {
			// 回调异常.
			defer func() {
				res.Finish = time.Now()

				// 捕获异常.
				if r := recover(); r != nil {
					res.Errors = append(res.Errors, fmt.Errorf("task panic: %v", r))
					res.Returned = nil
				} else {
					suc = true
				}
			}()

			// 执行过程.
			res.Returned = callback(ctx)
			return
		}(o.ctx, o.handler)
	}

	// 结果转发.
	if o.finish != nil {
		// 计算用时.
		res.Duration = res.Finish.Sub(res.Begin).Milliseconds()
		res.Delay = res.Begin.Sub(res.Created).Milliseconds()

		// 转发过程.
		func(ctx context.Context, callback TaskFinish, result TaskResult) {
			defer func() {
				if r := recover(); r != nil {
					log.Panicfc(ctx, "task panic: %v", r)
				}
			}()

			callback(ctx, result)
		}(o.ctx, o.finish, res)
	}

	// 释放实例.
	// 任务执行完成后自动释放.
	o.after()
	taskPool.Put(o)
	return
}

func (o *task) SetContext(ctx context.Context) Task {
	o.ctx = ctx
	return o
}

func (o *task) SetFinish(finish TaskFinish) Task {
	o.finish = finish
	return o
}

func (o *task) SetHandler(handler TaskHandler) Task {
	o.handler = handler
	return o
}

// /////////////////////////////////////////////////////////////
// Pool instance operations
// /////////////////////////////////////////////////////////////

func (o *task) after() {
	o.ctx = nil
	o.finish = nil
	o.handler = nil
}

func (o *task) before() *task {
	atomic.AddUint64(&o.acquires, 1)
	o.acquired = time.Now()
	return o
}

func (o *task) init() *task {
	o.id = atomic.AddUint64(&taskId, 1)
	return o
}
