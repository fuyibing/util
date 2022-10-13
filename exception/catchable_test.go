// author: wsfuyibing <websearch@163.com>
// date: 2022-10-11

package exception

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	x := New()
	x.Ignore(
		func(ctx context.Context) (ignored bool) {
			return
		},
	).Try(
		func(ctx context.Context) (ignored bool) {
			panic("panic in try")
			return
		},
	).Catch(func(ctx context.Context, err interface{}) (ignored bool) {
		return
	}).Finally(func(ctx context.Context) (ignored bool) {
		return
	}).Panic(func(ctx context.Context, v interface{}) {
	}).Run(nil)
}
