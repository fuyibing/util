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
	// IgnoreHandler
	// process callback handler.
	//
	// If return value is true, ignore next handlers in heap.
	IgnoreHandler func(ctx context.Context) (ignored bool)

	// PanicHandler
	// process callback handler.
	//
	// Called when panic occurred at runtime, then ignore next
	// handlers in heap.
	PanicHandler func(ctx context.Context, v interface{})

	// Processor
	// run like os process.
	Processor interface {
		// Add
		// child processes into heap.
		Add(ps ...Processor) Processor

		// After
		// register after handlers into heap.
		After(cs ...IgnoreHandler) Processor

		// Before
		// register before handlers into heap.
		Before(cs ...IgnoreHandler) Processor

		// Callback
		// register main handlers into heap.
		Callback(cs ...IgnoreHandler) Processor

		// Del
		// child processes from heap.
		Del(ps ...Processor) Processor

		// Get
		// return child process from heap.
		Get(name string) (process Processor, exists bool)

		// GetParent
		// return parent process.
		GetParent() (process Processor)

		// Healthy
		// return process health status.
		//
		// Return true if process context built and cancelled
		// signal not received.
		Healthy() bool

		// Name
		// return process name.
		//
		//   return "my-process"
		Name() string

		// Panic
		// register panic handler on process.
		Panic(cp PanicHandler) Processor

		// Restart process.
		Restart()

		// Start process.
		//
		// Return error if started already.
		Start(ctx context.Context) error

		// Stop process.
		Stop()

		// Stopped
		// return process status is stopped already.
		Stopped() bool

		// UnbindWhenStopped
		// config process unbind switch.
		UnbindWhenStopped(b bool) Processor

		// Bind
		// parent process on this.
		bind(p Processor) Processor
	}

	processor struct {
		cancel context.CancelFunc
		ctx    context.Context

		mu            sync.RWMutex
		name          string
		running, redo bool

		children   map[string]Processor
		ca, cb, cc []IgnoreHandler
		cp         PanicHandler
		parent     Processor
		parentUws  bool
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
func (o *processor) After(cs ...IgnoreHandler) Processor              { o.ca = cs; return o }
func (o *processor) Before(cs ...IgnoreHandler) Processor             { o.cb = cs; return o }
func (o *processor) Callback(cs ...IgnoreHandler) Processor           { o.cc = cs; return o }
func (o *processor) Del(ps ...Processor) Processor                    { return o.del(ps) }
func (o *processor) Get(name string) (process Processor, exists bool) { return o.get(name) }
func (o *processor) GetParent() (process Processor)                   { return o.getParent() }
func (o *processor) Healthy() bool                                    { return o.healthy() }
func (o *processor) Name() string                                     { return o.name }
func (o *processor) Panic(cp PanicHandler) Processor                  { o.cp = cp; return o }
func (o *processor) Restart()                                         { o.restart() }
func (o *processor) Start(ctx context.Context) error                  { return o.start(ctx) }
func (o *processor) Stop()                                            { o.stop() }
func (o *processor) Stopped() bool                                    { return o.stopped() }
func (o *processor) UnbindWhenStopped(b bool) Processor               { o.parentUws = b; return o }

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

// Add
// child processes into heap.
func (o *processor) add(ps []Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Range processes
	// then add to heap.
	for _, p := range ps {
		// Ignore if exists.
		if _, ok := o.children[p.Name()]; ok {
			continue
		}

		// Update children.
		o.children[p.Name()] = p.bind(o)
	}
	return o
}

// Bind
// parent process on this.
func (o *processor) bind(p Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.parent = p
	return o
}

// Del
// child processes from heap.
func (o *processor) del(ps []Processor) Processor {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Range processes
	// then delete from heap if exists.
	for _, p := range ps {
		if _, ok := o.children[p.Name()]; ok {
			delete(o.children, p.Name())
		}
	}
	return o
}

// Get
// return child process.
func (o *processor) get(name string) (process Processor, exists bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	process, exists = o.children[name]
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

// Init
// process fields.
func (o *processor) init() *processor {
	o.children = make(map[string]Processor)
	o.mu = sync.RWMutex{}
	o.parentUws = false
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

		// Send cancel signal.
		o.cancel()
	}
}

// Start process.
func (o *processor) start(ctx context.Context) (err error) {
	o.mu.Lock()

	// Return repeat running error.
	if o.running {
		o.mu.Unlock()
		return fmt.Errorf("process '%s' running", o.name)
	}

	// Set process status as running
	o.running = true
	o.mu.Unlock()

	// Set process status as stopped.
	defer func() {
		// Delete from parent.
		if o.parentUws && o.parent != nil {
			o.parent.Del(o)
		}

		// Revert status.
		o.mu.Lock()
		defer o.mu.Unlock()
		o.initState()
	}()

	// Call before handlers.
	if ci, ce := o.doHandlers(ctx, o.cb); ci {
		return ce
	}

	// Call after handlers, override err result if error
	// returned by any after handler.
	defer func(c context.Context) {
		if _, ce := o.doHandlers(c, o.ca); ce != nil && err == nil {
			err = ce
			return
		}
	}(ctx)

	// Loop call main handlers until process stop signal
	// received.
	for {
		// Return
		// for parent context cancelled error.
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
			_, ce := o.doHandlers(c, o.cc)
			return ce
		}(o.ctx, o.cancel)

		// Stop children, block coroutine until child processes
		// stopped.
		o.doChildStopped()

		// Destroy process context.
		o.mu.Lock()
		o.ctx = nil
		o.cancel = nil
		o.mu.Unlock()
	}
}

// Stop process.
func (o *processor) stop() {
	if o.healthy() {
		o.mu.Lock()
		o.redo = false
		o.mu.Unlock()

		// Send cancel signal.
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

// /////////////////////////////////////////////////////////////
// lifetime method.
// /////////////////////////////////////////////////////////////

func (o *processor) doCatch(ctx context.Context) (c, ci bool, ce error) {
	if v := recover(); v != nil {
		c = true
		ci = true
		ce = fmt.Errorf("%v", v)

		// Call panic callback.
		if o.cp != nil {
			o.cp(ctx, v)
		}
	}
	return
}

func (o *processor) doChildStart(ctx context.Context) {
	for _, child := range func() map[string]Processor {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return o.children
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
		return o.children
	}() {
		// Stopped already.
		if child.Stopped() {
			continue
		}

		// Send stop signal if health.
		if child.Healthy() {
			child.Stopped()
		}

		time.Sleep(time.Millisecond)
		return o.doChildStopped()
	}
	return true
}

func (o *processor) doHandlers(ctx context.Context, handlers []IgnoreHandler) (ignored bool, err error) {
	// Override result
	// if panic occurred.
	defer func() {
		if c, ci, ce := o.doCatch(ctx); c {
			err = ce
			ignored = ci
		}
	}()

	// Range handlers, break false returned by any handler
	for _, handler := range handlers {
		if ignored = handler(ctx); ignored {
			break
		}
	}

	return
}
