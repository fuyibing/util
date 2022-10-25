// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

var (
	Coder CoderManager
)

type (
	CoderManager interface {
		Integer(code Code) int
		SetPlus(plus int) CoderManager
	}

	coder struct {
		plus int
	}
)

func (o *coder) Integer(code Code) int {
	return o.plus + code.Int()
}

func (o *coder) SetPlus(plus int) CoderManager {
	o.plus = plus
	return o
}

func (o *coder) init() *coder {
	return o
}
