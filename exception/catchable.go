// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

package exception

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/fuyibing/util/v2/caller"
)

var (
	catchableId       uint64
	catchablePool     *sync.Pool
	catchableUnActive = fmt.Errorf("unactive catchable")
)

type (
	// Catchable
	// 可捕获接口.
	Catchable interface {
		// Catch
		// 注册捕获回调.
		Catch(cs ...caller.CatchCaller) Catchable

		// Finally
		// 注册最终回调.
		Finally(cs ...caller.FinallyCaller) Catchable

		// Identify
		// 获取实例ID.
		Identify() (id, acquires uint64)

		// Ignore
		// 注册忽略回调.
		Ignore(ci ...caller.IgnoreCaller) Catchable

		// Panic
		// 注册异常回调.
		Panic(cp caller.PanicCaller) Catchable

		// Run
		// 执行实例.
		Run(ctx context.Context) error

		// Try
		// 注册尝试回调.
		Try(cs ...caller.TryCaller) Catchable
	}

	// 可捕获结构体.
	catchable struct {
		acquires, id uint64
		active       bool
		mu           sync.RWMutex

		cc []caller.CatchCaller   // 可捕获回调列表
		cf []caller.FinallyCaller // 可最终回调列表
		ci []caller.IgnoreCaller  // 可忽略回调列表
		cp caller.PanicCaller     // 运行异常回调
		ct []caller.TryCaller     // 可尝试回调列表
	}
)

// New
// 创建实例.
//
// 创建实例时从池中获取.
func New() Catchable {
	return catchablePool.
		Get().(*catchable).
		before()
}

// Catch
// 注册捕获回调.
func (o *catchable) Catch(cs ...caller.CatchCaller) Catchable {
	o.cc = cs
	return o
}

// Finally
// 注册最终回调.
func (o *catchable) Finally(cs ...caller.FinallyCaller) Catchable {
	o.cf = cs
	return o
}

// Ignore
// 注册忽略回调.
func (o *catchable) Ignore(cs ...caller.IgnoreCaller) Catchable {
	o.ci = cs
	return o
}

// Panic
// 注册异常回调.
func (o *catchable) Panic(c caller.PanicCaller) Catchable {
	o.cp = c
	return o
}

// Try
// 注册尝试回调.
func (o *catchable) Try(cs ...caller.TryCaller) Catchable {
	o.ct = cs
	return o
}

// Identify
// 获取实例ID.
func (o *catchable) Identify() (id, acquires uint64) {
	return o.id, o.acquires
}

// Run
// 执行实例.
func (o *catchable) Run(ctx context.Context) (err error) {
	var ignored = false

	// 1. 返回错误.
	if func() bool {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return o.active == false
	}() {
		return catchableUnActive
	}

	// 2. 监听结束.
	//    捕获过程异常并清理数据, 最后释放实例回池.
	defer func() {
		if r := recover(); r != nil && o.cp != nil {
			o.cp(ctx, r)
		}
		o.after()
		catchablePool.Put(o)
	}()

	// 3. 前置执行.
	//    遍历可忽略回调, 任一回调返回 true 时退出, 因前置致退出时也会
	//    忽略已注册的 TryCaller/CatchCaller/FinallyCaller 回调.
	if len(o.ci) > 0 {
		for _, c := range o.ci {
			if ignored, err = o.doIgnore(ctx, c); ignored {
				return
			}
		}
	}

	// 4. 尝试回调.
	if len(o.ct) > 0 {
		// 4.1 触发尝试回调.
		for _, c := range o.ct {
			if ignored, err = o.doTry(ctx, c); ignored {
				break
			}
		}

		// 4.2 触发异常回调.
		if err != nil && len(o.cc) > 0 {
			for _, c := range o.cc {
				if o.doCatch(ctx, c, err) {
					break
				}
			}
		}

		// 4.3 触发最终回调.
		if len(o.cf) > 0 {
			for _, c := range o.cf {
				if o.doFinally(ctx, c) {
					break
				}
			}
		}
	}

	// 5. 完成
	return
}

// /////////////////////////////////////////////////////////////
// Callback callers execution.
// /////////////////////////////////////////////////////////////

func (o *catchable) doCatch(ctx context.Context, callback caller.CatchCaller, err error) (ignored bool) {
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

func (o *catchable) doFinally(ctx context.Context, callback caller.FinallyCaller) (ignored bool) {
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

func (o *catchable) doIgnore(ctx context.Context, callback caller.IgnoreCaller) (ignored bool, err error) {
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

func (o *catchable) doTry(ctx context.Context, callback caller.TryCaller) (ignored bool, err error) {
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
// Instance initializer.
// /////////////////////////////////////////////////////////////

func (o *catchable) after() *catchable {
	o.cc = nil
	o.cf = nil
	o.ci = nil
	o.cp = nil
	o.ct = nil
	o.mu.Lock()
	o.active = false
	o.mu.Unlock()
	return o
}

func (o *catchable) before() *catchable {
	atomic.AddUint64(&o.acquires, 1)
	o.mu.Lock()
	o.active = true
	o.mu.Unlock()
	return o
}

func (o *catchable) init() *catchable {
	o.id = atomic.AddUint64(&catchableId, 1)
	o.mu = sync.RWMutex{}
	return o
}
