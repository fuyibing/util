// author: wsfuyibing <websearch@163.com>
// date: 2023-02-20

package process

import (
	"context"
	"testing"
	"time"
)

var (
	rootCtx, rootCanceler = context.WithCancel(context.Background())
)

func TestNew(t *testing.T) {
	// Cancel by root context after 1 second..
	go func() {
		time.Sleep(time.Second)
		t.Logf("------------- root cancel")
		rootCanceler()
	}()

	m := (&my{t: t, name: "p1"}).init()
	if err := m.processor.Start(rootCtx); err != nil {
		t.Logf("%s stopped: %v", m.processor.Name(), err)
	}
}

func TestProcessor_Add(t *testing.T) {
	var m = (&my{t: t, name: "p1"}).init()

	// Cancel by root context after 3 second..
	go func() {

		for _, name := range []string{"c1", "c2", "c3"} {
			time.Sleep(time.Second)

			if child, exists := m.processor.Get(name); exists {
				t.Logf("---- [child=%s] restart", name)
				child.Restart()
			} else {
				t.Logf("---- [child=%s] not found", name)
			}
		}

		time.Sleep(time.Second * 2)
		t.Logf("---- ---- [root] restart")
		m.processor.Restart()

		time.Sleep(time.Second * 3)
		t.Logf("---- ---- [root] cancel context")
		rootCanceler()
	}()

	m = (&my{t: t, name: "p1"}).init()
	m.processor.Add(
		(&my{t: t, name: "c1"}).init().processor,
		(&my{t: t, name: "c2"}).init().processor.UnbindWhenStopped(true),
		(&my{t: t, name: "c3"}).init().processor,
	)

	if err := m.processor.Start(rootCtx); err != nil {
		t.Logf("%s stopped: %v", m.processor.Name(), err)
	}
}

// My struct for processor.

type my struct {
	name      string
	processor Processor
	t         *testing.T
}

func (o *my) init() *my {
	o.processor = New(o.name).
		After(o.onAfter).
		Before(o.onBefore).
		Callback(o.onCall).
		Panic(o.onPanic)
	return o
}

func (o *my) onAfter(_ context.Context) (ignored bool) {
	o.t.Logf("%s after", o.processor.Name())
	return
}

func (o *my) onBefore(_ context.Context) (ignored bool) {
	o.t.Logf("%s before", o.processor.Name())
	return
}

func (o *my) onCall(ctx context.Context) (ignored bool) {
	o.t.Logf("%s call", o.processor.Name())

	for {
		select {
		case <-ctx.Done():
			o.t.Logf("%s call: %v", o.processor.Name(), ctx.Err())
			return
		}
	}
}

func (o *my) onPanic(_ context.Context, v interface{}) {
	o.t.Logf("%s panic: %v", o.processor.Name(), v)
}
