// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

package exception

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	catchableId   uint64
	catchablePool *sync.Pool
)

type (
	// Catchable
	// 可捕获接口.
	Catchable interface {
		// Catch
		// 注册捕获回调.
		Catch(cs ...FuncCatch) Catchable

		// Finally
		// 注册最终回调.
		Finally(cs ...FuncFinally) Catchable

		// Ignore
		// 注册忽略回调.
		Ignore(ci ...FuncIgnore) Catchable

		// Panic
		// 注册异常回调.
		Panic(cp FuncPanic) Catchable

		// Run
		// 执行实例.
		Run(ctx context.Context) error

		// Try
		// 注册尝试回调.
		Try(cs ...FuncTry) Catchable
	}

	catchable struct {
		acquires, id uint64

		cc []FuncCatch   // 可捕获回调列表
		cf []FuncFinally // 可最终回调列表
		ci []FuncIgnore  // 可忽略回调列表
		cp FuncPanic     // 运行异常回调
		ct []FuncTry     // 可尝试回调列表
	}
)

// New
// 创建实例.
//
// 从池中获取实例, 当 Run() 被调用且执行完成后自动释放回池中.
func New() Catchable {
	return catchablePool.
		Get().(*catchable).
		before()
}

// Catch
// 注册捕获回调.
func (o *catchable) Catch(cs ...FuncCatch) Catchable {
	o.cc = cs
	return o
}

// Finally
// 注册最终回调.
func (o *catchable) Finally(cs ...FuncFinally) Catchable {
	o.cf = cs
	return o
}

// Ignore
// 注册忽略回调.
func (o *catchable) Ignore(cs ...FuncIgnore) Catchable {
	o.ci = cs
	return o
}

// Panic
// 注册异常回调.
func (o *catchable) Panic(c FuncPanic) Catchable {
	o.cp = c
	return o
}

// Try
// 注册尝试回调.
func (o *catchable) Try(cs ...FuncTry) Catchable {
	o.ct = cs
	return o
}

// Run
// 执行实例.
func (o *catchable) Run(ctx context.Context) (err error) {
	var ignored = false

	// 1. 监听结束.
	//    捕获过程异常并清理数据, 最后释放实例回池.
	defer func() {
		if r := recover(); r != nil && o.cp != nil {
			o.cp(ctx, r)
		}
		o.after()
		catchablePool.Put(o)
	}()

	// 2. 前置执行.
	//    遍历可忽略回调, 任一回调返回 true 时退出, 因前置致退出时也会
	//    忽略已注册的 FuncTry/FuncCatch/FuncFinally 回调.
	if len(o.ci) > 0 {
		for _, c := range o.ci {
			if ignored, err = o.doIgnore(ctx, c); ignored {
				return
			}
		}
	}

	// 3. 尝试回调.
	if len(o.ct) > 0 {
		// 3.1 触发尝试回调.
		for _, c := range o.ct {
			if ignored, err = o.doTry(ctx, c); ignored {
				break
			}
		}

		// 3.2 触发异常回调.
		if err != nil && len(o.cc) > 0 {
			for _, c := range o.cc {
				if o.doCatch(ctx, c, err) {
					break
				}
			}
		}

		// 3.3 触发最终回调.
		if len(o.cf) > 0 {
			for _, c := range o.cf {
				if o.doFinally(ctx, c) {
					break
				}
			}
		}
	}

	return nil
}

// /////////////////////////////////////////////////////////////
// Callbacks executor.
// /////////////////////////////////////////////////////////////

func (o *catchable) doCatch(ctx context.Context, callback FuncCatch, err error) (ignored bool) {
	defer func() {
		if r := recover(); r != nil {
			if o.cp != nil {
				o.cp(ctx, r)
			}
			ignored = true
		}
	}()
	ignored = callback(ctx, err)
	return
}

func (o *catchable) doFinally(ctx context.Context, callback FuncFinally) (ignored bool) {
	defer func() {
		if r := recover(); r != nil {
			if o.cp != nil {
				o.cp(ctx, r)
			}
			ignored = true
		}
	}()
	ignored = callback(ctx)
	return
}

func (o *catchable) doIgnore(ctx context.Context, callback FuncIgnore) (ignored bool, err error) {
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

func (o *catchable) doTry(ctx context.Context, callback FuncTry) (ignored bool, err error) {
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

// /////////////////////////////////////////////////////////////
// Not export methods.
// /////////////////////////////////////////////////////////////

func (o *catchable) after() *catchable {
	o.cc = nil
	o.cf = nil
	o.ci = nil
	o.cp = nil
	o.ct = nil
	return o
}

func (o *catchable) before() *catchable {
	atomic.AddUint64(&o.acquires, 1)
	return o
}

func (o *catchable) init() *catchable {
	o.id = atomic.AddUint64(&catchableId, 1)
	return o
}
