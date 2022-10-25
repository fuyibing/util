// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

func Data(v interface{}) *Result {
	return With.Data(v)
}

func Error(err error) *Result {
	return With.Error(err)
}

func ErrorCode(err error, code Code) *Result {
	return With.ErrorCode(err, code)
}

func Paging(v interface{}, total int64, limit, page int) *Result {
	return With.Paging(v, total, limit, page)
}

func Succeed() *Result {
	return With.Succeed()
}
