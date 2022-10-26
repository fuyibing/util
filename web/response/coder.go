// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

var (
	Coder CodeManager
)

type (
	CodeManager interface {
		// Integer
		// 返回错误码.
		//
		// 当错误在[100-999)间有效.
		Integer(code int) int

		// SetPlus
		// 设置错误码基数.
		//
		// 仅对非HTTP状态码[100-999)有效.
		SetPlus(plus int) CodeManager
	}

	code struct {
		plus int
	}
)

// Integer
// 返回错误码.
func (o *code) Integer(n int) int {
	if n >= 100 && n < 1000 {
		return n
	}
	return o.plus + n
}

// SetPlus
// 设置错误码基数.
func (o *code) SetPlus(n int) CodeManager {
	o.plus = n
	return o
}

func (o *code) init() *code {
	return o
}
