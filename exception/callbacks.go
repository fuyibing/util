// author: wsfuyibing <websearch@163.com>
// date: 2022-10-12

package exception

import "context"

type (
	// FuncCatch
	// 捕获回调.
	//
	// 此回调在Try/Catch块中触发, 当业务执行过程中出现panic时自动执行[0-1]次, 若
	// Try回调无panic或Catch列表中前一个返回true时执行0次, 反之执行1次.
	FuncCatch func(ctx context.Context, err interface{}) (ignored bool)

	// FuncFinally
	// 最终回调.
	//
	// 此回调在Try/Catch块后触发, 最多执行[0-1]次, 若注册了 FuncIgnore 且返回 true
	// 或Finally列表中前一个返回 true 时执行0次, 反之执行1次.
	FuncFinally func(ctx context.Context) (ignored bool)

	// FuncIgnore
	// 忽略回调.
	//
	// 此回调在Try/Catch块前触发, 最多执行[0-1]次, 若前一个返回 true 时执行0次,
	// 反之执行1次.
	FuncIgnore func(ctx context.Context) (ignored bool)

	// FuncPanic
	// 异常回调.
	//
	// 当业务执行过程中出现panic时自动触发, 通常用于记录日志, 在Try/Catch块中此回
	// 调可能执行[0-N]次.
	FuncPanic func(ctx context.Context, v interface{})

	// FuncTry
	// 偿试回调.
	//
	// 此回调在Try/Catch块中触发, 最多执行[0-1]次, 若注册了 FuncIgnore 且返回 true
	// 或Try列表中前一个返回 true 时执行0次, 反之执行1次.
	FuncTry func(ctx context.Context) (ignored bool)
)
