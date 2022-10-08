// author: wsfuyibing <websearch@163.com>
// date: 2020-01-01

package util

import (
    "context"
    "fmt"
    "sync"
)

type (
    // Processor
    // 执行器接口.
    //
    // 模拟多进程服务中的 Worker 进程（实际运行在协程 Goroutine 中）, 生命周期从 Start() 开始,
    // 到 context.CancelFunc 触发后结束.
    //
    // 需安全退出(确定指定任务执行完成)时, 可通过 After() 注册回调.
    Processor interface {
        // After
        // 注册后置回调.
        //
        // 执行器退出前触发, 至少执行 1 次.
        After(cs ...SkipCaller) Processor

        // Before
        // 注册前置回调.
        //
        // 启动执行器前触发, 至少执行 1 次.
        Before(cs ...SkipCaller) Processor

        // Callee
        // 注册过程回调.
        //
        // 执行器主体执行过程, 执行 0-n 次, 若通过 Before() 注册的回调返回了忽略(Skipped)则
        // 不执行(执行0次), 首次启动或每次重启时各执行 1 次.
        Callee(cs ...SkipCaller) Processor

        // Healthy
        // 健康状态检查.
        //
        // 启动成功并且 context 未收到 Cancel 信号.
        Healthy() bool

        // Panic
        // 注册异常回调.
        //
        // 执行过程中, 每次出现 panic, 此回调触发1次.
        Panic(cp PanicCaller) Processor

        // Restart
        // 重启执行器.
        Restart()

        // Running
        // 执行器状态.
        //
        // 包含已启动或退出中状态.
        Running() bool

        // Start
        // 启动执行器.
        Start(ctx context.Context) error

        // Stop
        // 退出执行器.
        Stop()
    }

    processor struct {
        cancel           context.CancelFunc
        ctx              context.Context
        mu               sync.RWMutex
        name             string
        running, restart bool

        ca, cb, cc []SkipCaller
        cp         PanicCaller
    }
)

func NewProcessor(name string) Processor {
    return (&processor{name: name}).init()
}

// After
// 注册后置回调.
func (o *processor) After(cs ...SkipCaller) Processor {
    o.ca = append(o.ca, cs...)
    return o
}

// Before
// 注册前置回调.
func (o *processor) Before(cs ...SkipCaller) Processor {
    o.cb = append(o.cb, cs...)
    return o
}

// Callee
// 注册过程回调.
func (o *processor) Callee(cs ...SkipCaller) Processor {
    o.cc = append(o.cc, cs...)
    return o
}

// Healthy
// 健康状态检查.
func (o *processor) Healthy() bool {
    o.mu.RLock()
    defer o.mu.RUnlock()
    return o.ctx != nil && o.ctx.Err() == nil
}

// Panic
// 注册异常回调.
func (o *processor) Panic(cp PanicCaller) Processor {
    o.cp = cp
    return o
}

// Restart
// 重启执行器.
func (o *processor) Restart() {
    o.mu.Lock()
    defer o.mu.Unlock()
    if o.ctx != nil && o.ctx.Err() == nil {
        o.restart = true
        o.cancel()
    }
}

// Running
// 执行器状态.
func (o *processor) Running() bool {
    o.mu.RLock()
    defer o.mu.RUnlock()
    return o.running
}

// Start
// 启动执行器.
func (o *processor) Start(ctx context.Context) error {
    o.mu.Lock()
    if o.running {
        o.mu.Unlock()
        return fmt.Errorf("started already: %s", o.name)
    }
    o.running = true
    o.mu.Unlock()

    defer func() {
        if r := recover(); r != nil && o.cp != nil {
            o.cp(ctx, r)
        }

        o.onStop(ctx)

        o.mu.Lock()
        o.running = false
        o.mu.Unlock()
    }()

    if o.onStart(ctx) == nil {
        o.onProgress(ctx)
    }
    return nil
}

// Stop
// 退出执行器.
func (o *processor) Stop() {
    o.mu.Lock()
    defer o.mu.Unlock()
    if o.ctx != nil && o.ctx.Err() == nil {
        o.restart = false
        o.cancel()
    }
}

// /////////////////////////////////////////////////////////////
// Constructor
// /////////////////////////////////////////////////////////////

func (o *processor) init() *processor {
    o.mu = sync.RWMutex{}
    o.ca = make([]SkipCaller, 0)
    o.cb = make([]SkipCaller, 0)
    o.cc = make([]SkipCaller, 0)
    return o
}

// /////////////////////////////////////////////////////////////
// Events
// /////////////////////////////////////////////////////////////

func (o *processor) onProgress(ctx context.Context) {
    if ctx == nil || ctx.Err() != nil {
        return
    }

    o.mu.Lock()
    o.ctx, o.cancel = context.WithCancel(ctx)
    o.mu.Unlock()

    defer func() {
        if r := recover(); r != nil && o.cp != nil {
            o.cp(ctx, r)
        }

        if o.ctx.Err() != nil {
            o.cancel()
        }
        o.ctx = nil
        o.cancel = nil

        if o.restart {
            o.restart = false
            o.onProgress(ctx)
        }
    }()

    for _, caller := range o.cc {
        if caller(o.ctx) {
            break
        }
    }
}

func (o *processor) onStart(ctx context.Context) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("%v", r)

            if o.cp != nil {
                o.cp(ctx, r)
            }
        }
    }()

    for _, call := range o.cb {
        if call(ctx) {
            break
        }
    }
    return
}

func (o *processor) onStop(ctx context.Context) {
    defer func() {
        if r := recover(); r != nil && o.cp != nil {
            o.cp(ctx, r)
        }
    }()

    for _, call := range o.ca {
        if call(ctx) {
            break
        }
    }
}
