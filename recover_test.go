// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package util

import (
	"github.com/fuyibing/log/v3/trace"
	"testing"
)

func TestRecover(t *testing.T) {
	Recover(callback1, callback2, callback3)
}

func TestRecoverWithContext(t *testing.T) {
	ctx := trace.New()
	RecoverWithContext(ctx, callback1, callback2, callback3)
}

func callback1() {
	println("callback 1")
}

func callback2() {
	println("callback 2")
	panic("callback2")
}

func callback3() {
	println("callback 3")
}
