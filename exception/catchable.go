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
		Catch(cs ...Catch) Catchable

		// Finally
		// 注册最终回调.
		Finally(cs ...Finally) Catchable

		// Ignore
		// 注册忽略回调.
		Ignore(ci ...Ignore) Catchable

		// Panic
		// 注册异常回调.
		Panic(cp Panic) Catchable

		// Run
		// 执行实例.
		Run(ctx context.Context) error

		// Try
		// 注册尝试回调.
		Try(cs ...Try) Catchable
	}

	catchable struct {
		acquires, id uint64

		cc []Catch   // 可捕获回调列表
		cf []Finally // 可最终回调列表
		ci []Ignore  // 可忽略回调列表
		cp Panic     // 运行异常回调
		ct []Try     // 可尝试回调列表
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
func (o *catchable) Catch(cs ...Catch) Catchable {
	o.cc = cs
	return o
}

// Finally
// 注册最终回调.
func (o *catchable) Finally(cs ...Finally) Catchable {
	o.cf = cs
	return o
}

// Ignore
// 注册忽略回调.
func (o *catchable) Ignore(cs ...Ignore) Catchable {
	o.ci = cs
	return o
}

// Panic
// 注册异常回调.
func (o *catchable) Panic(c Panic) Catchable {
	o.cp = c
	return o
}

// Try
// 注册尝试回调.
func (o *catchable) Try(cs ...Try) Catchable {
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
	//    忽略已注册的 Try/Catch/Finally 回调.
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

func (o *catchable) doCatch(ctx context.Context, callback Catch, err error) (ignored bool) {
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

func (o *catchable) doFinally(ctx context.Context, callback Finally) (ignored bool) {
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

func (o *catchable) doIgnore(ctx context.Context, callback Ignore) (ignored bool, err error) {
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

func (o *catchable) doTry(ctx context.Context, callback Try) (ignored bool, err error) {
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
