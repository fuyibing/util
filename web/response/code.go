// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

var (
	Code *CodeManager
)

const (
	UndefinedCode = 1
)

type (
	CodeManager struct {
		plus int
	}
)

// Integer
// 返回错误码.
func (o *CodeManager) Integer(n int) int {
	if n >= 100 && n < 1000 {
		return n
	}
	return o.plus + n
}

// SetPlus
// 设置错误码基数.
func (o *CodeManager) SetPlus(n int) *CodeManager {
	o.plus = n
	return o
}

func (o *CodeManager) init() *CodeManager {
	return o
}
