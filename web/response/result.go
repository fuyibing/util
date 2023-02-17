// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

package response

import (
	"encoding/json"
)

type (
	Result struct {
		Data     interface{} `json:"data" label:"Data"`
		DataType Type        `json:"data_type" label:"Data type"`
		Errno    int         `json:"errno" label:"Error number"`
		Error    string      `json:"error" label:"Error reason"`
	}

	Type string
)

const (
	TypeData   Type = "OBJECT"
	TypeError  Type = "ERROR"
	TypeList   Type = "LIST"
	TypePaging Type = "PAGING"

	ResultFieldForBody  = "body"
	ResultNameForPaging = "paging"
)

func NewResult(dt Type) *Result {
	return &Result{DataType: dt}
}

func (o *Result) Json() string {
	buf, _ := json.Marshal(o)
	return string(buf)
}
