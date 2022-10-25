// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

var (
	With WithManager
)

type (
	WithManager interface {
		Data(v interface{}) *Result
		Error(err error) *Result
		ErrorCode(err error, code Code) *Result
		List(v interface{}) *Result
		Paging(v interface{}, total int64, limit, page int) *Result
		Succeed() *Result
	}

	with struct {
	}
)

var (
	defaultOnData = make(map[string]interface{})
	defaultOnList = make([]interface{}, 0)
)

func (o *with) Data(data interface{}) *Result {
	r := NewResult(TypeData)

	if data == nil {
		r.Data = defaultOnData
	} else {
		r.Data = data
	}

	return r
}

func (o *with) Error(err error) *Result {
	return o.ErrorCode(err, UndefinedCode)
}

func (o *with) ErrorCode(err error, code Code) *Result {
	r := NewResult(TypeError)
	r.Errno = Coder.Integer(code)
	r.Error = err.Error()
	return r
}

func (o *with) List(v interface{}) *Result {
	r := NewResult(TypeList)

	if v == nil {
		r.Data = defaultOnList
	} else {
		r.Data = v
	}

	return r
}

func (o *with) Paging(v interface{}, total int64, limit, page int) *Result {
	r := NewResult(TypePaging)
	return r
}

func (o *with) Succeed() *Result {
	return o.Data(nil)
}

// 构建实例.
func (o *with) init() *with {
	return o
}
