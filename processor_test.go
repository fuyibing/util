// author: wsfuyibing <websearch@163.com>
// date: 2020-01-01

package util

import (
    "context"
    "testing"
    "time"
)

type processorTester struct {
}

func (o *processorTester) a1(ctx context.Context) bool {
    println("after 1")
    return false
}

func (o *processorTester) a2(ctx context.Context) bool {
    println("after 2")
    return true
}

func (o *processorTester) a3(ctx context.Context) bool {
    println("after 3")
    return false
}

func (o *processorTester) b1(ctx context.Context) bool {
    println("before 1")
    panic("before 1")
    return false
}

func (o *processorTester) b2(ctx context.Context) bool {
    println("before 2")
    return false
}

func (o *processorTester) b3(ctx context.Context) bool {
    println("before 3")
    return false
}

func (o *processorTester) c1(ctx context.Context) bool {
    println("callee 1")
    return false
}

func (o *processorTester) c2(ctx context.Context) bool {
    for {
        select {
        case <-ctx.Done():
            println("callee 2:", ctx.Err().Error())
            return false
        }
    }
}

func (o *processorTester) c3(ctx context.Context) bool {
    println("callee 3")
    return true
}

func (o *processorTester) c4(ctx context.Context) bool {
    println("callee 4")
    return false
}

func TestProcessor(t *testing.T) {
    var (
        ctx, cancel = context.WithCancel(context.TODO())
        err         error
        p           = NewProcessor("test-processor")
        x           = &processorTester{}
    )

    go func() {
        for i := 0; i < 3; i++ {
            time.Sleep(time.Second * 1)
            p.Restart()
        }

        time.Sleep(time.Second * 1)
        cancel()

        t.Errorf("error: %v.", err)
    }()

    err = p.After(x.a1, x.a2, x.a3).Before(x.b1, x.b2, x.b3).Callee(x.c1, x.c2, x.c3, x.c4).
        Panic(func(_ context.Context, v interface{}) {
            t.Errorf("panic: %v.", v)
        }).Start(ctx)
}
