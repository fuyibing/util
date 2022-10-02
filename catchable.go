// author: wsfuyibing <websearch@163.com>
// date: 2022-10-02

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
    // interface for try/catch block runner.
    Catchable interface {
        Before(cs ...SkipCaller) Catchable
        Catch(cs ...CatchCaller) Catchable
        Finally(cs ...FinallyCaller) Catchable
        Identify() (id, acquires uint64)
        Logger(c PanicCaller) Catchable
        Run(ctx context.Context) error
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
// create catchable instance.
func TryCatch() Catchable {
    return catchablePool.Get().(*catchable).before()
}

// Before
// register skip-doCaller.
func (o *catchable) Before(cs ...SkipCaller) Catchable {
    o.beforeCallers = append(o.beforeCallers, cs...)
    return o
}

// Catch
// register catch-doCaller.
func (o *catchable) Catch(cs ...CatchCaller) Catchable {
    o.catchCallers = append(o.catchCallers, cs...)
    return o
}

// Finally
// register finally-doCaller.
func (o *catchable) Finally(cs ...FinallyCaller) Catchable {
    o.finallyCallers = append(o.finallyCallers, cs...)
    return o
}

// Identify
// return instance id of pool and acquired count from pool.
func (o *catchable) Identify() (id, acquires uint64) {
    return o.id, o.acquires
}

// Logger
// register panic doCaller.
func (o *catchable) Logger(c PanicCaller) Catchable {
    o.panicCaller = c
    return o
}

// Run
// catchable progress and return error if panic occurred otherwise
// nil returned.
func (o *catchable) Run(ctx context.Context) error {
    return o.run(ctx)
}

// Try
// register try-doCaller.
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

    // Try/Cache before doCaller.
    skip := false
    if skip, err = o.runBefore(ctx); skip {
        return
    }

    // Run try doCaller of try/catch.
    if err = o.runTry(ctx); err != nil {
        // Run catch doCaller of try/catch.
        o.runCatch(ctx, err)
    }

    // Run finally doCaller of try/catch.
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
    // and break if true returned in any doCaller.
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
    // and break if true returned in any doCaller.
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
    // and break if true returned in any doCaller.
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
    // and break if true returned in any doCaller.
    for _, caller := range o.tryCallers {
        if caller(ctx) {
            break
        }
    }
    return
}
