// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

// Package exception
//
//   ctx := log.NewContext()
//
//   obj := exception.New().Ignore(
//        func(ctx context.Context) (ignored bool) {
//            return
//        },
//    ).Try(
//        func(ctx context.Context) (ignored bool) {
//            panic("panic in try")
//            return
//        },
//    ).Catch(func(ctx context.Context, err interface{}) (ignored bool) {
//        return
//    }).Finally(func(ctx context.Context) (ignored bool) {
//        return
//    }).Panic(func(ctx context.Context, v interface{}) {
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
