// author: wsfuyibing <websearch@163.com>
// date: 2020-01-01

package util

import "context"

type (
    // CatchCaller
    // 捕获回调.
    CatchCaller func(ctx context.Context, e error) (skipped bool)

    // FinallyCaller
    // 最终回调.
    //
    // 若前置回调
    FinallyCaller func(ctx context.Context) (skipped bool)

    // TryCaller
    // 偿试回调.
    TryCaller func(ctx context.Context) (skipped bool)

    // PanicCaller
    // 异常回调.
    //
    // 仅在运行阶段出现 Panic 时触发, 每次出现 Panic 必触发一次, 通常用于记录日志.
    PanicCaller func(ctx context.Context, v interface{})

    // SkipCaller
    // 可忽略回调.
    //
    // 当返回 true 时, 表示忽略后续的回调.
    SkipCaller func(ctx context.Context) (skipped bool)
)
