// author: wsfuyibing <websearch@163.com>
// date: 2022-10-02

package util

import "context"

type (
    // CatchCaller
    // register as catch-caller of Try/Catch.
    CatchCaller func(ctx context.Context, e error) (skipped bool)

    // FinallyCaller
    // register as finally-caller of Try/Catch.
    FinallyCaller func(ctx context.Context) (skipped bool)

    // TryCaller
    // register as try-caller of Try/Catch.
    TryCaller func(ctx context.Context) (skipped bool)

    // PanicCaller
    // register as panic-caller when panic occurred in any caller.
    PanicCaller func(ctx context.Context, v interface{})

    // SkipCaller
    // register as skip-caller, skip next callers if true returned.
    SkipCaller func(ctx context.Context) (skipped bool)
)
