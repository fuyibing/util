// author: wsfuyibing <websearch@163.com>
// date: 2022-10-23

package response

type (
	// Result
	// 返回结果.
	Result struct {
		Data     interface{} `json:"data" label:"结果数据"`
		DataType Type        `json:"data_type" label:"数据类型"`
		Errno    int         `json:"errno" label:"错误编码"`
		Error    string      `json:"error" label:"错误原因"`
	}

	// Type
	// 结果类型.
	Type string
)

const (
	TypeData   Type = "OBJECT"
	TypeError  Type = "ERROR"
	TypeList   Type = "LIST"
	TypePaging Type = "PAGING"
)

// NewResult
// 创建返回结果.
func NewResult(dt Type) *Result {
	return &Result{DataType: dt}
}
