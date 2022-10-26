// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

import (
	"fmt"
	"net/http"
	"testing"
)

func ExampleData() {
	v := With.Succeed()
	println("data: ", v.String())
}

func ExampleError() {
	err := fmt.Errorf("example error")
	res := With.Error(err)
	println("result = ", res.String())
}

func ExampleErrorCode() {
	Coder.SetPlus(3000)
	err := fmt.Errorf("example error")
	code := http.StatusForbidden
	res := With.ErrorCode(err, code)
	println("result = ", res.String())
}

func ExampleList() {
}

func ExamplePaging() {}

func ExampleSucceed() {
}

func TestData(t *testing.T)      { ExampleData() }
func TestError(t *testing.T)     { ExampleError() }
func TestErrorCode(t *testing.T) { ExampleErrorCode() }
func TestList(t *testing.T)      { ExampleList() }
func TestPaging(t *testing.T)    { ExamplePaging() }
func TestSucceed(t *testing.T)   { ExampleSucceed() }
