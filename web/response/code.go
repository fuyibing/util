// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

type (
	// Code
	// 错误码.
	Code int
)

// 错误码枚举.
const (
	_ Code = iota

	UndefinedCode // 内部错误.
)

func (o Code) Int() int {
	return int(o)
}
