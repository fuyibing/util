// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

package exception

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
	// 1. 创建实例.
	//    a). 从池中获取
	//    b). 若池中实例不足则创建新实例
	catch := New()

	// 2. 注册回调/可忽略回调列表.
	//    a). 触发次数: [0-1]
	//    b). 回调列表中任一返回 true 时结束执行.
	//    c). 回调列表中任一位置出现 panic 时结束执行.
	catch.Ignore(
		// 回调 #1
		func(ctx context.Context) (ignored bool) {
			return
		},
		// 回调 #2
		func(ctx context.Context) (ignored bool) {
			return
		},
	)

	// 3. 注册回调/尝试回调列表.
	//    a). 触发次数: [0-1]
	//    b). 未注册忽略回调或全部返回 false.
	//    c). 回调列表任一回调返回 true, 忽略后续尝试回调.
	//    d). 回调列表中任一位置出现 panic 时忽略后续尝试回调并触发捕获回调.
	catch.Try(
		// 回调 #1
		func(ctx context.Context) (ignored bool) {
			time.Sleep(time.Millisecond * 3)
			return
		},
		// 回调 #2
		func(ctx context.Context) (ignored bool) {
			return
		},
	)

	// 4. 注册回调/捕获回调列表.
	//    a). 触发次数: [0-1]
	//    b). 未注册忽略回调或全部返回 false.
	//    c). 尝试回调列表任一位置出现过panic异常.
	catch.Catch(
		// 回调 #1
		func(ctx context.Context, err interface{}) (ignored bool) {
			return
		},
		// 回调 #2
		func(ctx context.Context, err interface{}) (ignored bool) {
			return
		},
	)

	// 5. 注册回调/最终回调列表.
	//    a). 触发次数: [0-1]
	//    b). 尝试回调触发过.
	catch.Finally(
		// 回调 #1
		func(ctx context.Context) (ignored bool) {
			return
		},
		// 回调 #1
		func(ctx context.Context) (ignored bool) {
			return
		},
	)

	// 6. 注册回调/异常回调.
	//    a). 触发次数: [0-N]
	//    b). 整个生命周期中, 任一位置出现 panic 时都会触发本回调
	catch.Panic(
		func(ctx context.Context, v interface{}) {
		},
	)

	// 7. 获取实例ID.
	//    a). 返回值id为实例唯一编号
	//    b). 返回值acquires为实例从池中取出次数
	id, acquires := catch.Identify()
	fmt.Printf("info: id=%d, acquires=%d.\n", id, acquires)

	// 8. 实例执行过程.
	//    a). 实例执行过程.
	//    b). 执行完成自动释放回池.
	ctx := context.TODO()
	if err := catch.Run(ctx); err != nil {
		fmt.Printf("run: %v.\n", err)
	}
	fmt.Printf("run: completed.\n")
}
