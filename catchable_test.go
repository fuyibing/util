// author: wsfuyibing <websearch@163.com>
// date: 2020-01-01

package util

import (
    "context"
    "sync"
    "testing"
    "time"
)

func TestCatchable(t *testing.T) {
    if err := TryCatch().Before(
        func(ctx context.Context) (skip bool) {
            t.Logf("before-1")
            return
        },
        func(ctx context.Context) (skip bool) {
            t.Logf("before-2")
            return
        },
    ).Try(
        func(ctx context.Context) (skip bool) {
            t.Logf("try-1")
            panic("try-1")
            return
        },
        func(ctx context.Context) (skip bool) {
            t.Logf("try-2")
            return
        },
    ).Catch(
        func(ctx context.Context, e error) (skipped bool) {
            t.Logf("catch-1")
            panic("catch-1")
            return
        },
        func(ctx context.Context, e error) (skip bool) {
            t.Logf("catch-2")
            return
        },
    ).Finally(
        func(ctx context.Context) (skip bool) {
            t.Logf("finally-1")
            return
        },
        func(ctx context.Context) (skip bool) {
            t.Logf("finally-2")
            return
        },
    ).Panic(
        func(ctx context.Context, v interface{}) {
            t.Logf("panic: %v", v)
        },
    ).Run(context.TODO()); err != nil {
        t.Errorf("try/catch: %v.", err)
        return
    }
    t.Logf("try/catch: completed.")
}

func TestCatchableGoroutine(t *testing.T) {
    ctx := context.TODO()

    for i := 0; i < 10; i++ {
        x := TryCatch()
        m, n := x.Identify()
        _ = x.Run(ctx)
        t.Logf("cacheable: %d -> %d.", m, n)
    }
}

func TestCatchableGoroutines(t *testing.T) {
    ctx := context.TODO()
    wait := sync.WaitGroup{}
    for i := 0; i < 10; i++ {
        wait.Add(1)
        go func() {
            defer wait.Done()
            x := TryCatch()
            m, n := x.Identify()
            _ = x.Run(ctx)
            t.Logf("cacheable: %d -> %d.", m, n)
            time.Sleep(time.Millisecond * 5)
        }()
    }
    wait.Wait()
}
