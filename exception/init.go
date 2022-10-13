// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

// Package exception
// 异常捕获.
package exception

import "sync"

func init() {
	new(sync.Once).Do(func() {
		catchablePool = &sync.Pool{
			New: func() interface{} { return (&catchable{}).init() },
		}
	})
}
