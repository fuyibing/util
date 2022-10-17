// author: wsfuyibing <websearch@163.com>
// date: 2022-10-13

package process

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ExampleNew()
}

func ExampleNew() {

	var (
		ctx, cancel = context.WithCancel(context.TODO())
		err         error
		proc        = New("example")
	)

	// 1. 延时处理.
	//    a). 延时 3 秒, 重启进程.
	//    b). 延时 10 秒, 退出进程.
	//    c). 延时 15 秒, 退出所有进程.
	// 延时5秒后强制退出.
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 3)
			fmt.Printf("-------------- restart #%d ---------\n", i)
			proc.Restart()
		}

		time.Sleep(time.Second * 5)
		fmt.Println("-------------- stop ---------")
		proc.Stop()

		time.Sleep(time.Second * 5)
		fmt.Println("-------------- stop all ---------")
		cancel()
	}()

	// 2. 注册回调/前置过程.
	proc.After(
		// 回调 #1
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("after #1.\n")
			return
		},
		// 回调 #2
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("after #2.\n")
			return
		},
	)

	// 3. 注册回调/后置过程.
	proc.Before(
		// 回调 #1
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("before #1.\n")
			return
		},
		// 回调 #2
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("before #2.\n")
			return
		},
	)

	// 4. 注册回调/过程处理.
	proc.Callback(
		// 回调 1
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("callback #1.\n")
			return
		},

		// 回调 2
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("callback #2.\n")
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("callback #2: end.\n")
					return
				}
			}
		},

		// 回调 3
		func(ctx context.Context) (ignored bool) {
			fmt.Printf("callback #3.\n")
			panic("callback 3")
			return
		},
	)

	proc.Panic(func(ctx context.Context, v interface{}) {
		fmt.Printf("panic: %v.\n", v)
	})

	// 5. 启动进程.
	fmt.Println("-------------- start ---------")
	if err = proc.Start(ctx); err != nil {
		fmt.Printf("process: %v.\n", err)
		return
	}
	fmt.Printf("process: completed.\n")
}
