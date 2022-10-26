// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

import (
	"fmt"
	"net/http"
	"testing"
)

func ExampleWithManager_Data() {
	data := map[string]interface{}{
		"id":  1,
		"key": "value",
	}
	v := With.Data(data)
	println("result: ", v.Json())
}
func ExampleWithManager_Error() {
	err := fmt.Errorf("error message")
	res := With.Error(err)
	println("result:", res.Json())
}
func ExampleWithManager_ErrorCode() {
	err := fmt.Errorf("forbidden")
	res := With.ErrorCode(err, http.StatusForbidden)
	println("result:", res.Json())
}
func ExampleWithManager_List() {
	v := []map[string]interface{}{
		{"id": 1, "key": "value-1"},
		{"id": 2, "key": "value-2"},
	}
	res := With.List(v)
	println("result:", res.Json())
}
func ExampleWithManager_Paging() {
	v := []map[string]interface{}{
		{"id": 1, "key": "value-1"},
		{"id": 2, "key": "value-2"},
	}
	res := With.Paging(v, 2, 10, 1)
	println("result:", res.Json())
}
func ExampleWithManager_Success() {
	res := With.Success()
	println("result:", res.Json())
}

func TestWithManager_Data(t *testing.T)      { ExampleWithManager_Data() }
func TestWithManager_Error(t *testing.T)     { ExampleWithManager_Error() }
func TestWithManager_ErrorCode(t *testing.T) { ExampleWithManager_ErrorCode() }
func TestWithManager_List(t *testing.T)      { ExampleWithManager_List() }
func TestWithManager_Paging(t *testing.T)    { ExampleWithManager_Paging() }
func TestWithManager_Success(t *testing.T)   { ExampleWithManager_Success() }
