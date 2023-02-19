// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

// Package process
// run like os process.
package process

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	ErrRunningAlready = fmt.Errorf("running already")
)

type (
	Handler      func(ctx context.Context) (ignored bool)
	PanicHandler func(ctx context.Context, v interface{})

	// Processor
	// run like os process.
	Processor interface {
		// Add
		// child processor to this.
		//
		// Start child processors when this processor started.
		Add(ps ...Processor) Processor

		// After
		// register handlers when processor stop.
		//
		// Called frequency: 1.
		After(cs ...Handler) Processor

		// Before
		// register handlers when processor start.
		//
		// Called frequency: 1.
		Before(cs ...Handler) Processor

		// Callback
		// register main handlers when processor start or restart.
		//
		// Called frequency: N.
		Callback(cs ...Handler) Processor

		// Del
		// remove child process.
		Del(ps ...Processor) Processor

		// Healthy
		// return processor status.
		Healthy() bool

		// Name
		// return processor name.
		Name() string

		// Panic
		// register panic handler.
		//
		// Called frequency: N.
		Panic(c PanicHandler) Processor

		// RemoveFromParent
		// remove from parent if enabled..
		//
		// Default: false
		RemoveFromParent(rm bool) Processor

		// Restart
		// stop handlers which registered by Callback method,
		// then start again.
		Restart()

		// Start
		// send start signal to this.
		Start(ctx context.Context) error

		// Stop
		// send stop signal to this.
		Stop()

		// Stopped
		// return processor stopped status.
		Stopped() bool
	}

	processor struct {
		cancel context.CancelFunc
		ctx    context.Context

		mu               sync.RWMutex
		name             string
		parent           Processor
		parentRemove     bool
		running, restart bool

		children   map[string]Processor
		ca, cb, cc []Handler
		cp         PanicHandler
	}
)

// New
// create and return processor interface.
func New(name string) Processor {
	return (&processor{
		children: make(map[string]Processor, 0),
		name:     name,
	}).init()
}

func (o *processor) Add(ps ...Processor) Processor      { return o.add(ps...) }
func (o *processor) After(cs ...Handler) Processor      { o.ca = cs; return o }
func (o *processor) Del(ps ...Processor) Processor      { return o.del(ps...) }
func (o *processor) Before(cs ...Handler) Processor     { o.cb = cs; return o }
func (o *processor) Callback(cs ...Handler) Processor   { o.cc = cs; return o }
func (o *processor) Name() string                       { return o.name }
func (o *processor) Panic(c PanicHandler) Processor     { o.cp = c; return o }
func (o *processor) RemoveFromParent(rm bool) Processor { o.parentRemove = rm; return o }

func (o *processor) Healthy() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.ctx != nil && o.ctx.Err() == nil
}

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

func (o *processor) Start(ctx context.Context) error {
	o.mu.Lock()

	// Return error
	// if process was started already.
	if o.running {
		o.mu.Unlock()
		return ErrRunningAlready
	}

	// Update running status.
	o.running = true
	o.mu.Unlock()

	// Auto called
	// before process end.
	defer func() {
		// Send panic error
		// if necessary .
		if r := recover(); r != nil && o.cp != nil {
			o.cp(ctx, r)
		}

		// Cancel context
		// if necessary.
		if o.ctx != nil && o.ctx.Err() == nil {
			o.cancel()
		}

		// Iterate registered after handlers.
		for _, c := range o.ca {
			if ci, _ := o.doHandler(ctx, c); ci {
				break
			}
		}

		// Revert running status.
		o.mu.Lock()
		o.running = false
		o.mu.Unlock()

		// Remove from parent.
		if o.parentRemove && o.parent != nil {
			o.parent.Del(o)
		}
	}()

	// Call registered before handlers. Quit next handlers
	// if ignored returned.
	for _, c := range o.cb {
		if ci, ce := o.doHandler(ctx, c); ci {
			return ce
		}
	}

	// Call registered callback handlers.
	for {
		// Return
		// if not restart signal.
		if func() bool {
			o.mu.RLock()
			defer o.mu.RUnlock()
			return o.restart == false
		}() {
			return nil
		}

		// Revert
		// restart status.
		//
		// Very Important.
		o.mu.Lock()
		o.restart = false
		o.mu.Unlock()

		// Return
		// if parent context cancelled.
		if ctx == nil || ctx.Err() != nil {
			return nil
		}

		// Registered main handlers executors.
		func() {
			o.mu.Lock()
			o.ctx, o.cancel = context.WithCancel(ctx)
			o.mu.Unlock()

			// Auto called
			// when main handlers execute finish.
			defer func() {
				// Force cancel context.
				if o.ctx.Err() == nil {
					o.cancel()
				}

				// Block and wait
				// child processors stopped.
				o.childWait()

				// Revert process context.
				o.mu.Lock()
				o.ctx = nil
				o.cancel = nil
				o.mu.Unlock()
			}()

			// Start child processors.
			o.childStart(o.ctx)

			// Iterate registered main handlers. Break if
			// ignored return.
			for _, c := range o.cc {
				if ci, _ := o.doHandler(o.ctx, c); ci {
					break
				}
			}
		}()
	}
}

func (o *processor) Stop() {
	o.mu.Lock()
	if o.ctx != nil && o.ctx.Err() == nil {
		o.restart = false
		o.mu.Unlock()
		o.cancel()
		return
	}
	o.mu.Unlock()
}

func (o *processor) Stopped() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.running == false
}

// /////////////////////////////////////////////////////////////
// Callback callers execution.
// /////////////////////////////////////////////////////////////

func (o *processor) add(ps ...Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, p := range ps {
		k := p.Name()
		if _, ok := o.children[k]; ok {
			continue
		}

		p.(*processor).parent = o
		o.children[k] = p
	}

	return o
}

func (o *processor) del(ps ...Processor) Processor {
	for _, p := range ps {
		go o.delProcessor(p)
	}
	return o
}

func (o *processor) delProcessor(p Processor) {
	for {
		if p.Stopped() {
			o.mu.Lock()
			delete(o.children, p.Name())
			o.mu.Unlock()
			break
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func (o *processor) doHandler(ctx context.Context, handler Handler) (ignored bool, err error) {
	// Catch handler panic
	// when end.
	defer func() {
		r := recover()

		// Success called.
		if r == nil {
			return
		}

		// Send panic error.
		if o.cp != nil {
			o.cp(ctx, r)
		}

		// Override handler result.
		ignored = true
		err = fmt.Errorf("%v", r)
	}()

	// Call registered handler.
	ignored = handler(ctx)
	return
}

func (o *processor) childStart(ctx context.Context) {
	for _, c := range o.children {
		go func(p Processor) {
			_ = p.Start(ctx)
		}(c)
	}
}

func (o *processor) childWait() {
	for _, c := range o.children {
		if !c.Stopped() {
			time.Sleep(time.Millisecond * 10)
			o.childWait()
			return
		}
	}
}

// /////////////////////////////////////////////////////////////
// Constructor method.
// /////////////////////////////////////////////////////////////

func (o *processor) init() *processor {
	o.mu = sync.RWMutex{}
	o.restart = true
	o.running = false
	return o
}
