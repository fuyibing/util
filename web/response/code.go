// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

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
// return target code.
func (o *CodeManager) Integer(n int) int {
	if n >= 100 && n < 1000 {
		return n
	}
	return o.plus + n
}

// SetPlus
// config base integer for plus.
func (o *CodeManager) SetPlus(n int) *CodeManager {
	o.plus = n
	return o
}

func (o *CodeManager) init() *CodeManager {
	return o
}
