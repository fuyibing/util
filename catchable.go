// author: wsfuyibing <websearch@163.com>
// date: 2020-01-01

package util

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
)

var (
    catchableId   uint64
    catchablePool sync.Pool
)

type (
    // Catchable
    // 捕获异常操作接口.
    //
    // 模拟 Try/Catch 代码块, 在此调用中使用注册回调的方式运行.
    Catchable interface {
        // Before
        // 注册前置回调.
        //
        // 前置回调列表中, 任一回调返回 true, 则忽略已注册的 TryCaller, CatchCaller,
        // FinallyCaller 回调(不执行).
        Before(cs ...SkipCaller) Catchable

        // Catch
        // 注册捕获回调.
        //
        // 当执行 TryCaller 时, 任一回调出现 panic 时, 退出 TryCaller 并立即
        // 触发 CatchCaller. 若整个过程未出现过 panic, 则不执行 CatchCaller.
        Catch(cs ...CatchCaller) Catchable

        // Finally
        // 注册最终回调.
        //
        // 当由 Before() 注册的 SkipCaller 执行完成后(无 true 返回), FinallyCaller
        // 必触发。
        Finally(cs ...FinallyCaller) Catchable

        // Identify
        // 返回ID.
        //
        // 返回此实例(来自池)的ID号, 从池中取出次数. 例如 acquires=10, 表示此实例第
        // 10 次从池中取出(复用10次).
        Identify() (id, acquires uint64)

        // Panic
        // 注册异常回调.
        Panic(c PanicCaller) Catchable

        // Run
        // 执行过程.
        //
        // 当返回值为非 nil 时, 表示执行过程中出现 panic.
        Run(ctx context.Context) error

        // Try
        // 注册尝试回调.
        Try(cs ...TryCaller) Catchable
    }

    catchable struct {
        acquires, id uint64

        beforeCallers  []SkipCaller
        catchCallers   []CatchCaller
        finallyCallers []FinallyCaller
        panicCaller    PanicCaller
        tryCallers     []TryCaller
    }
)

// TryCatch
// 创建异常操作实例.
func TryCatch() Catchable {
    return catchablePool.Get().(*catchable).before()
}

// Before
// 注册前置回调.
func (o *catchable) Before(cs ...SkipCaller) Catchable {
    o.beforeCallers = append(o.beforeCallers, cs...)
    return o
}

// Catch
// 注册捕获回调.
func (o *catchable) Catch(cs ...CatchCaller) Catchable {
    o.catchCallers = append(o.catchCallers, cs...)
    return o
}

// Finally
// 注册最终回调.
func (o *catchable) Finally(cs ...FinallyCaller) Catchable {
    o.finallyCallers = append(o.finallyCallers, cs...)
    return o
}

// Identify
// 返回ID.
func (o *catchable) Identify() (id, acquires uint64) {
    return o.id, o.acquires
}

// Panic
// 注册异常回调.
func (o *catchable) Panic(c PanicCaller) Catchable {
    o.panicCaller = c
    return o
}

// Run
// 执行过程.
//
// 当返回值为非 nil 时, 表示执行过程中出现 panic.
func (o *catchable) Run(ctx context.Context) error {
    return o.run(ctx)
}

// Try
// 注册尝试回调.
func (o *catchable) Try(cs ...TryCaller) Catchable {
    o.tryCallers = append(o.tryCallers, cs...)
    return o
}

// Invisible called
// when release to pool.
func (o *catchable) after() *catchable {
    o.panicCaller = nil

    o.beforeCallers = nil
    o.catchCallers = nil
    o.finallyCallers = nil
    o.tryCallers = nil

    catchablePool.Put(o)
    return o
}

// Invisible called
// when acquired from pool.
func (o *catchable) before() *catchable {
    atomic.AddUint64(&o.acquires, 1)

    o.beforeCallers = make([]SkipCaller, 0)
    o.catchCallers = make([]CatchCaller, 0)
    o.finallyCallers = make([]FinallyCaller, 0)
    o.tryCallers = make([]TryCaller, 0)
    return o
}

// Invisible called
// when instance in pool not enough.
func (o *catchable) init() *catchable {
    o.id = atomic.AddUint64(&catchableId, 1)
    return o
}

// Run progress.
func (o *catchable) run(ctx context.Context) (err error) {
    // Release to pool
    // when progress end.
    defer o.after()

    // Try/Cache before caller.
    skip := false
    if skip, err = o.runBefore(ctx); skip {
        return
    }

    // Run try caller of try/catch.
    if err = o.runTry(ctx); err != nil {
        // Run catch caller of try/catch.
        o.runCatch(ctx, err)
    }

    // Run finally caller of try/catch.
    o.runFinally(ctx)
    return
}

// Run SkipCaller
// which registered by Before() method.
func (o *catchable) runBefore(ctx context.Context) (skip bool, err error) {
    // Catch panic.
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("%v", r)
            skip = true

            if o.panicCaller != nil {
                o.panicCaller(ctx, r)
            }
        }
    }()

    // Run skip-callers
    // and break if true returned in any caller.
    for _, caller := range o.beforeCallers {
        if skip = caller(ctx); skip {
            break
        }
    }
    return
}

// Run CatchCaller
// which registered by Catch() method.
func (o *catchable) runCatch(ctx context.Context, v error) {
    // Catch panic.
    defer func() {
        if r := recover(); r != nil {
            if o.panicCaller != nil {
                o.panicCaller(ctx, r)
            }
        }
    }()

    // Run catch-callers
    // and break if true returned in any caller.
    for _, caller := range o.catchCallers {
        if caller(ctx, v) {
            break
        }
    }
}

// Run FinallyCaller
// which registered by Finally() method.
func (o *catchable) runFinally(ctx context.Context) {
    // Catch panic.
    defer func() {
        if r := recover(); r != nil {
            if o.panicCaller != nil {
                o.panicCaller(ctx, r)
            }
        }
    }()

    // Run finally-callers
    // and break if true returned in any caller.
    for _, caller := range o.finallyCallers {
        if caller(ctx) {
            break
        }
    }
    return
}

// Run TryCaller
// which registered by Try() method.
func (o *catchable) runTry(ctx context.Context) (err error) {
    // Catch panic.
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("%v", r)

            if o.panicCaller != nil {
                o.panicCaller(ctx, r)
            }
        }
    }()

    // Run try-callers
    // and break if true returned in any caller.
    for _, caller := range o.tryCallers {
        if caller(ctx) {
            break
        }
    }
    return
}
