// author: wsfuyibing <websearch@163.com>
// date: 2022-10-26

package response

import (
	"testing"
)

func ExampleNewPaginator() {
	var (
		total       int64 = 123
		limit, page       = 10, 3
	)

	p := NewPaginator(total, limit, page)
	println("result: ", p.Json())
}

func TestNewPaginator(t *testing.T) {
	ExampleNewPaginator()

	// for i := 0; i <= 101; i++ {
	// 	p := NewPaginator(int64(i), 10, 2)
	// 	t.Logf("p: %d->%+v", i, p.TotalPages)
	// }
}
