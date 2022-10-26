// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package util

import (
	"context"
	"github.com/fuyibing/log/v3"
)

// Recover
// 捕获运行异常.
func Recover(callbacks ...func()) {
	RecoverWithContext(nil, callbacks...)
}

// RecoverWithContext
// 捕获运行异常.
func RecoverWithContext(ctx context.Context, callbacks ...func()) {
	// 1. 捕获异常.
	defer func() {
		if r := recover(); r != nil {
			log.Panicfc(ctx, "runtime panic: %v", r)
		}
	}()

	// 2. 遍历回调.
	for _, callback := range callbacks {
		callback()
	}
}
