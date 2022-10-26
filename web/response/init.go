// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

// Package response
// 返回结果处理.
package response

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Code = (&CodeManager{}).init()
		With = (&WithManager{}).init()
	})
}
