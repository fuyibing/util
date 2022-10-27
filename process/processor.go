// author: wsfuyibing <websearch@163.com>
// date: 2022-10-13

package process

import (
	"context"
	"fmt"
	"sync"

	"github.com/fuyibing/util/v2/caller"
)

type (
	// Processor
	// 模拟进程接口.
	Processor interface {
		// After
		// 后置回调.
		//
		// 执行次数: 1 次.
		// 若回调列表中任一回调返回 true 或出现 panic 时忽略后续回调.
		After(cs ...caller.IgnoreCaller) Processor

		// Before
		// 前置回调.
		//
		// 执行次数: 1 次.
		// 若回调列表中任一回调返回 true 或出现 panic 时忽略后续回调并退出进程, 且
		// 不执行 Callback() 注册过的进程回调.
		Before(cs ...caller.IgnoreCaller) Processor

		// Callback
		// 进程回调.
		//
		// 执行次数: [0-n] 次
		Callback(cs ...caller.ProcessCaller) Processor

		// Healthy
		// 返回健康状态.
		//
		// 进程已启动且未收到退出信息时返回 true, 反之返回 false.
		Healthy() bool

		// Panic
		// 捕获异常.
		//
		// 执行次数: [0-n] 次
		// 进程的生命周期中, 任一位置出现 panic 时, 都会触发过回调.
		Panic(c caller.PanicCaller) Processor

		// Restart
		// 重启进程.
		//
		// 仅对健康的进程有效.
		Restart()

		// Start
		// 启动进程.
		//
		// 仅对未启动或已完全退出的进程有效, 若进程已启动或退出中时返回错误.
		Start(ctx context.Context) error

		// Stop
		// 退出进程.
		//
		// 仅对健康的进程有效.
		Stop()

		// Stopped
		// 返回退出状态.
		//
		// 若进程从未启动或退出完成时返回 true, 反之返回 false.
		Stopped() bool
	}

	// 模拟进程结构体.
	processor struct {
		cancel           context.CancelFunc
		ctx              context.Context
		mu               sync.RWMutex
		name             string
		running, restart bool

		ca, cb []caller.IgnoreCaller
		cc     []caller.ProcessCaller
		cp     caller.PanicCaller
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
	o.ca = cs
	return o
}

// Before
// 注册前置回调.
func (o *processor) Before(cs ...caller.IgnoreCaller) Processor {
	o.cb = cs
	return o
}

// Callback
// 注册进程回调.
func (o *processor) Callback(cs ...caller.ProcessCaller) Processor {
	o.cc = cs
	return o
}

// Healthy
// 返回健康状态.
func (o *processor) Healthy() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.ctx != nil && o.ctx.Err() == nil
}

// Panic
// 注册异常回调.
func (o *processor) Panic(c caller.PanicCaller) Processor {
	o.cp = c
	return o
}

// Start
// 启动进程.
func (o *processor) Start(ctx context.Context) error {
	o.mu.Lock()

	// 1. 重复启动.
	if o.running {
		o.mu.Unlock()
		return errRunningAlready
	}

	// 2. 开始启动.
	o.running = true
	o.mu.Unlock()

	// 3. 监听结束.
	defer func() {
		// 3.1 捕获异常.
		if r := recover(); r != nil && o.cp != nil {
			o.cp(ctx, r)
		}

		// 3.2 退上下文.
		if o.ctx != nil && o.ctx.Err() == nil {
			o.cancel()
		}

		// 3.3 后置执行.
		for _, c := range o.ca {
			if ci, _ := o.doIgnore(ctx, c); ci {
				break
			}
		}

		// 3.4 结束进程.
		o.mu.Lock()
		o.running = false
		o.mu.Unlock()
	}()

	// 4. 前置执行.
	for _, c := range o.cb {
		if ci, ce := o.doIgnore(ctx, c); ci {
			return ce
		}
	}

	// 5. 执行过程.
	for {
		// 5.1 退出启动.
		if func() bool {
			o.mu.RLock()
			defer o.mu.RUnlock()
			return o.restart == false
		}() {
			return nil
		}

		// 5.2 上游退出.
		if ctx == nil || ctx.Err() != nil {
			return nil
		}

		// 5.3 启动过程.
		o.mu.Lock()
		o.restart = false
		o.mu.Unlock()

		func() {
			o.mu.Lock()
			o.ctx, o.cancel = context.WithCancel(ctx)
			o.mu.Unlock()

			defer func() {
				o.mu.Lock()
				o.ctx = nil
				o.cancel = nil
				o.mu.Unlock()
			}()

			for _, c := range o.cc {
				if o.doProcess(o.ctx, c) {
					break
				}
			}
		}()
	}
}

// Stop
// 退出进程.
func (o *processor) Stop() {
	o.mu.Lock()
	if o.ctx != nil && o.ctx.Err() == nil {
		o.restart = false
		o.mu.Unlock()
		o.cancel()
		return
	}
	o.mu.RUnlock()
}

// Stopped
// 返回退出状态.
func (o *processor) Stopped() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.running == false
}

// Restart
// 重启进程.
func (o *processor) Restart() {
	o.mu.Lock()
	if o.ctx != nil && o.ctx.Err() == nil {
		o.restart = true
		o.mu.Unlock()
		o.cancel()
		return
	}
	o.mu.RUnlock()
}

// /////////////////////////////////////////////////////////////
// Callback callers execution.
// /////////////////////////////////////////////////////////////

func (o *processor) doIgnore(ctx context.Context, callback caller.IgnoreCaller) (ignored bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if o.cp != nil {
				o.cp(ctx, r)
			}
			ignored = true
			err = fmt.Errorf("%v", r)
		}
	}()
	ignored = callback(ctx)
	return
}

func (o *processor) doProcess(ctx context.Context, callback caller.ProcessCaller) (ignored bool) {
	defer func() {
		if r := recover(); r != nil && o.cp != nil {
			ignored = true
			o.cp(ctx, r)
		}
	}()
	ignored = callback(ctx)
	return
}

// /////////////////////////////////////////////////////////////
// Instance initializer.
// /////////////////////////////////////////////////////////////

func (o *processor) init() *processor {
	o.mu = sync.RWMutex{}
	o.restart = true
	o.running = false
	return o
}
