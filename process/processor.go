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

type (
	// Event
	// auto called in process lifetime.
	//
	// Process will ignore next events in heap if return value is
	// false, otherwise call next events by registered order.
	Event func(ctx context.Context) (ignored bool)

	// PanicEvent
	// auto called if panic occurred in event.
	PanicEvent func(ctx context.Context, v interface{})

	// Processor
	// run like os process.
	Processor interface {
		// Add
		// subprocesses into process.
		Add(ps ...Processor) Processor

		// After
		// register after events.
		After(es ...Event) Processor

		// Before
		// register before events.
		Before(es ...Event) Processor

		// Callback
		// register main events.
		Callback(es ...Event) Processor

		// Del
		// subprocesses from process.
		Del(ps ...Processor) Processor

		// Get
		// return subprocess of process.
		Get(name string) (process Processor, exists bool)

		// GetParent
		// return parent process, return nil if this role
		// is not a subprocess.
		GetParent() (process Processor)

		// Healthy
		// return health status.
		//
		// Return true if process context built and cancelled
		// signal never received.
		Healthy() bool

		// Name
		// return process name.
		//
		//   return "my-process"
		Name() string

		// Panic
		// register panic event.
		Panic(cp PanicEvent) Processor

		// Restart process.
		Restart()

		// Start process.
		//
		// Return error if started already or is starting or is
		// stopping or is restarting.
		Start(ctx context.Context) error

		// StartChild start subprocess.
		StartChild(name string) error

		// Stop process.
		Stop()

		// Stopped
		// return process status is stopped already, return true if
		// never start.
		Stopped() bool

		// Unbind
		// call parent process delete child.
		Unbind() Processor

		// UnbindWhenStopped
		// config process unbind type.
		//
		// If true set, notify parent process remove subprocess
		// when stopped.
		UnbindWhenStopped(b bool) Processor

		// Bind
		// parent event on this.
		bind(p Processor) Processor
	}

	processor struct {
		cancel context.CancelFunc
		ctx    context.Context

		mu            sync.RWMutex
		name          string
		running, redo bool

		ae, be, ce     []Event
		pe             PanicEvent
		parent         Processor
		subprocesses   map[string]Processor
		unbindWhenStop bool
	}
)

// New
// create and return processor interface.
//
//   proc := process.New("my-process")
//   proc.Start(ctx)
func New(name string) Processor {
	return (&processor{name: name}).
		init()
}

// /////////////////////////////////////////////////////////////
// Interface method.
// /////////////////////////////////////////////////////////////

func (o *processor) Add(ps ...Processor) Processor                    { return o.add(ps) }
func (o *processor) After(cs ...Event) Processor                      { o.ae = cs; return o }
func (o *processor) Before(cs ...Event) Processor                     { o.be = cs; return o }
func (o *processor) Callback(cs ...Event) Processor                   { o.ce = cs; return o }
func (o *processor) Del(ps ...Processor) Processor                    { return o.del(ps) }
func (o *processor) Get(name string) (process Processor, exists bool) { return o.get(name) }
func (o *processor) GetParent() (process Processor)                   { return o.getParent() }
func (o *processor) Healthy() bool                                    { return o.healthy() }
func (o *processor) Name() string                                     { return o.name }
func (o *processor) Panic(cp PanicEvent) Processor                    { o.pe = cp; return o }
func (o *processor) Restart()                                         { o.restart() }
func (o *processor) Start(ctx context.Context) error                  { return o.start(ctx) }
func (o *processor) StartChild(name string) error                     { return o.startChild(name) }
func (o *processor) Stop()                                            { o.stop() }
func (o *processor) Stopped() bool                                    { return o.stopped() }
func (o *processor) Unbind() Processor                                { return o.unbind() }
func (o *processor) UnbindWhenStopped(b bool) Processor               { o.unbindWhenStop = b; return o }

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

// Add
// subprocesses into process.
func (o *processor) add(ps []Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, p := range ps {
		if _, ok := o.subprocesses[p.Name()]; ok {
			continue
		}
		o.subprocesses[p.Name()] = p.bind(o)
	}
	return o
}

// Bind
// parent event on this.
func (o *processor) bind(p Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.parent = p
	return o
}

// Del
// subprocesses from process.
func (o *processor) del(ps []Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, p := range ps {
		if _, ok := o.subprocesses[p.Name()]; ok {
			delete(o.subprocesses, p.Name())
		}
	}
	return o
}

// Get
// return subprocess.
func (o *processor) get(name string) (process Processor, exists bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	process, exists = o.subprocesses[name]
	return
}

// GetParent
// return parent process.
func (o *processor) getParent() (process Processor) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.parent
}

// Healthy
// return process health status.
func (o *processor) healthy() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.ctx != nil && o.ctx.Err() == nil
}

func (o *processor) init() *processor {
	o.subprocesses = make(map[string]Processor)
	o.mu = sync.RWMutex{}
	o.unbindWhenStop = false
	o.initState()
	return o
}

func (o *processor) initState() {
	o.redo = true
	o.running = false
}

// Restart process.
func (o *processor) restart() {
	if o.healthy() {
		o.mu.Lock()
		o.redo = true
		o.mu.Unlock()
		o.cancel()
	}
}

// Start process.
func (o *processor) start(ctx context.Context) (err error) {
	o.mu.Lock()

	// Return repeat running error.
	if o.running {
		o.mu.Unlock()
		return fmt.Errorf("process '%s' was started already", o.name)
	}

	// Set process status as running
	o.running = true
	o.mu.Unlock()

	// Set process status as stopped.
	defer func() {
		// Delete from parent.
		if o.unbindWhenStop && o.parent != nil {
			o.parent.Del(o)
		}

		o.mu.Lock()
		o.initState()
		o.mu.Unlock()
	}()

	// Call before events.
	if ci, ce := o.doHandlers(ctx, o.be); ci {
		return ce
	}

	// Call after events, override result if error returned by
	// any event.
	defer func(c context.Context) {
		if _, ce := o.doHandlers(c, o.ae); ce != nil && err == nil {
			err = ce
			return
		}
	}(ctx)

	// Loop call main handlers until process stop signal
	// received.
	for {
		// Return
		// for parent context cancelled.
		if ctx == nil || ctx.Err() != nil {
			return
		}

		// Return
		// for stop signal received.
		if !func() (re bool) {
			o.mu.Lock()
			defer o.mu.Unlock()
			if re = o.redo; re {
				o.redo = false
			}
			return
		}() {
			return
		}

		// Build process context.
		o.mu.Lock()
		o.ctx, o.cancel = context.WithCancel(ctx)
		o.mu.Unlock()

		// Start children.
		o.doChildStart(o.ctx)

		// Call main handlers.
		err = func(c context.Context, cc context.CancelFunc) error {
			defer cc()
			_, ce := o.doHandlers(c, o.ce)
			return ce
		}(o.ctx, o.cancel)

		// Stop subprocesses, block coroutine until all
		// subprocesses stopped.
		o.doChildStopped()

		// Destroy process context.
		o.mu.Lock()
		o.ctx = nil
		o.cancel = nil
		o.mu.Unlock()
	}
}

// startChild start subprocess.
func (o *processor) startChild(name string) (err error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	p, exists := o.subprocesses[name]
	if !exists {
		err = fmt.Errorf("subprocess '%s' not found", name)
		return
	}

	if !p.Stopped() {
		err = fmt.Errorf("subprocess '%s' started already", name)
		return
	}

	go func() { _ = p.Start(o.ctx) }()
	return
}

// Stop process.
func (o *processor) stop() {
	if o.healthy() {
		o.mu.Lock()
		o.redo = false
		o.mu.Unlock()
		o.cancel()
	}
}

// Stopped
// return process stopped status.
func (o *processor) stopped() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return !o.running
}

func (o *processor) unbind() *processor {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.parent != nil {
		o.parent.Del(o)
	}
	return o
}

// /////////////////////////////////////////////////////////////
// lifetime method.
// /////////////////////////////////////////////////////////////

func (o *processor) doChildStart(ctx context.Context) {
	for _, child := range func() map[string]Processor {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return o.subprocesses
	}() {
		if child.Stopped() {
			go func(c context.Context, p Processor) {
				_ = p.Start(c)
			}(ctx, child)
		}
	}
}

func (o *processor) doChildStopped() (stopped bool) {
	for _, child := range func() map[string]Processor {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return o.subprocesses
	}() {
		if child.Stopped() {
			continue
		}

		if child.Healthy() {
			child.Stopped()
		}

		time.Sleep(time.Millisecond)
		return o.doChildStopped()
	}
	return true
}

func (o *processor) doHandlers(ctx context.Context, handlers []Event) (ignored bool, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			ignored = true

			if o.pe != nil {
				o.pe(ctx, v)
			}
		}
	}()

	for _, handler := range handlers {
		if ignored = handler(ctx); ignored {
			break
		}
	}

	return
}
