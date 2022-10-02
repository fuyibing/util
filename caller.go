// author: wsfuyibing <websearch@163.com>
// date: 2022-10-02

package util

import "context"

type (
    // CatchCaller
    // register as catch-doCaller of Try/Catch.
    CatchCaller func(ctx context.Context, e error) (skipped bool)

    // FinallyCaller
    // register as finally-doCaller of Try/Catch.
    FinallyCaller func(ctx context.Context) (skipped bool)

    // TryCaller
    // register as try-doCaller of Try/Catch.
    TryCaller func(ctx context.Context) (skipped bool)

    // PanicCaller
    // register as panic-doCaller when panic occurred in any doCaller.
    PanicCaller func(ctx context.Context, v interface{})

    // SkipCaller
    // register as skip-doCaller, skip next callers if true returned.
    SkipCaller func(ctx context.Context) (skipped bool)
)
