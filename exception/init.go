// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

// Package exception
//
//   ctx := log.NewContext()
//
//   obj := exception.New().FuncIgnore(
//        func(ctx context.Context) (ignored bool) {
//            return
//        },
//    ).FuncTry(
//        func(ctx context.Context) (ignored bool) {
//            panic("panic in try")
//            return
//        },
//    ).FuncCatch(func(ctx context.Context, err interface{}) (ignored bool) {
//        return
//    }).FuncFinally(func(ctx context.Context) (ignored bool) {
//        return
//    }).FuncPanic(func(ctx context.Context, v interface{}) {
//        println("panic occurred")
//    })
//   obj.Run(ctx)
package exception

import "sync"

func init() {
	new(sync.Once).Do(func() {
		catchablePool = &sync.Pool{
			New: func() interface{} { return (&catchable{}).init() },
		}
	})
}
