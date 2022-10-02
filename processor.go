// author: wsfuyibing <websearch@163.com>
// date: 2022-10-02

package util

import (
    "context"
    "fmt"
    "sync"
)

type (
    // Processor
    // interface for process manager.
    Processor interface {
        After(cs ...SkipCaller) Processor
        Before(cs ...SkipCaller) Processor
        Callee(cs ...SkipCaller) Processor

        Healthy() bool
        Panic(cp PanicCaller) Processor
        Restart()
        Running() bool

        Start(ctx context.Context) error
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
// register callers when processor finish.
func (o *processor) After(cs ...SkipCaller) Processor {
    o.ca = append(o.ca, cs...)
    return o
}

// Before
// register callers when processor begin.
func (o *processor) Before(cs ...SkipCaller) Processor {
    o.cb = append(o.cb, cs...)
    return o
}

// Callee
// register callers when processor progressing.
func (o *processor) Callee(cs ...SkipCaller) Processor {
    o.cc = append(o.cc, cs...)
    return o
}

// Healthy
// return processor health status.
func (o *processor) Healthy() bool {
    o.mu.RLock()
    defer o.mu.RUnlock()
    return o.ctx != nil && o.ctx.Err() == nil
}

// Panic
// register caller when panic occurred.
func (o *processor) Panic(cp PanicCaller) Processor {
    o.cp = cp
    return o
}

// Restart
// processor progress.
func (o *processor) Restart() {
    o.mu.Lock()
    defer o.mu.Unlock()
    if o.ctx != nil && o.ctx.Err() == nil {
        o.restart = true
        o.cancel()
    }
}

// Running
// return processor status not stop completed.
func (o *processor) Running() bool {
    o.mu.RLock()
    defer o.mu.RUnlock()
    return o.running
}

// Start
// processor progress with context.
func (o *processor) Start(ctx context.Context) error {
    // Progress status.
    o.mu.Lock()
    if o.running {
        o.mu.Unlock()
        return fmt.Errorf("started already: %s", o.name)
    }
    o.running = true
    o.mu.Unlock()

    // Trigger
    // when processor end or panic occurred.
    defer func() {
        // Catch panic.
        if r := recover(); r != nil && o.cp != nil {
            o.cp(ctx, r)
        }

        // Event
        // on processor stop.
        o.onStop(ctx)

        // Processor end.
        o.mu.Lock()
        o.running = false
        o.mu.Unlock()
    }()

    // Event
    // on start and progress.
    if o.onStart(ctx) == nil {
        o.onProgress(ctx)
    }
    return nil
}

// Stop
// processor progress.
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
    // Return
    // if parent context cancelled.
    if ctx == nil || ctx.Err() != nil {
        return
    }

    // Build
    // processor context.
    o.mu.Lock()
    o.ctx, o.cancel = context.WithCancel(ctx)
    o.mu.Unlock()

    // Triggered
    // when progress finish.
    defer func() {
        // Catch panic.
        if r := recover(); r != nil && o.cp != nil {
            o.cp(ctx, r)
        }

        // Cancel
        // processor context.
        if o.ctx.Err() != nil {
            o.cancel()
        }
        o.ctx = nil
        o.cancel = nil

        // Restart
        // progress if Restart() method called.
        if o.restart {
            o.restart = false
            o.onProgress(ctx)
        }
    }()

    // Iterate callers.
    for _, caller := range o.cc {
        if caller(o.ctx) {
            break
        }
    }
}

func (o *processor) onStart(ctx context.Context) (err error) {
    // Execute and catch panic
    // when event finish.
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("%v", r)

            if o.cp != nil {
                o.cp(ctx, r)
            }
        }
    }()

    // Call before callers.
    for _, call := range o.cb {
        if call(ctx) {
            break
        }
    }
    return
}

func (o *processor) onStop(ctx context.Context) {
    // Execute and catch panic
    // when event finish.
    defer func() {
        if r := recover(); r != nil && o.cp != nil {
            o.cp(ctx, r)
        }
    }()

    // Call before callers.
    for _, call := range o.ca {
        if call(ctx) {
            break
        }
    }
}
