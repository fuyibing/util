// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package response

var (
	// With
	// 使用结果类型.
	With *WithManager

	defaultOnData = make(map[string]interface{})
)

type (
	// WithManager
	// 使用类型管理.
	WithManager struct{}
)

// Data
// 返回基础数据.
//
// [Code]:
//   data := map[string]interface{}{
//       "id": 1,
//       "key": "value"
//   }
//   response.With.Data(data)
//
// [Return]
//   return {
//       "data": {
//           "id": 1,
//           "key": "value"
//       },
//       "data_type": "OBJECT",
//       "errno": 0,
//       "error": ""
//   }
func (o *WithManager) Data(data interface{}) *Result {
	r := NewResult(TypeData)
	r.Data = data
	return r
}

// Error
// 返回错误数据.
//
// [Code]:
//   err := fmt.Errorf("error message")
//   response.With.Error(err)
//
// [Return]
//   return {
//       "data": {},
//       "data_type": "ERROR",
//       "errno": 1,
//       "error": "error message"
//   }
func (o *WithManager) Error(err error) *Result {
	return o.ErrorCode(err, UndefinedCode)
}

// ErrorCode
// 返回错误数据.
//
// [Code]:
//   err := fmt.Errorf("forbidden")
//   response.With.ErrorCode(err, http.StatusForbidden)
//
// [Return]
//   return {
//       "data": {},
//       "data_type": "ERROR",
//       "errno": 1,
//       "error": "error message"
//   }
func (o *WithManager) ErrorCode(err error, code int) *Result {
	r := NewResult(TypeError)
	r.Data = defaultOnData
	r.Errno = Code.Integer(code)
	r.Error = err.Error()
	return r
}

// List
// 返回列表数据.
//
// [Code]:
//   v := []map[string]interface{}{
//       {"id": 1, "key": "value-1"},
//       {"id": 2, "key": "value-2"},
//   }
//   response.With.List(v)
//
// [Return]:
//   return {
//       "data": {
//           "body": [
//               {
//               	"id": 1,
//               	"key": "value-1",
//               },
//               {
//               	"id": 2,
//               	"key": "value-2",
//               }
//           ]
//       },
//       "data_type": "LIST",
//       "errno": 0,
//       "error": "",
//   }
func (o *WithManager) List(v interface{}) *Result {
	r := NewResult(TypeList)
	r.Data = map[string]interface{}{ResultFieldForBody: v}
	return r
}

// Paging
// 返回分页数据.
//
// [Code]:
//   v := []map[string]interface{}{
//       {"id": 1, "key": "value-1"},
//       {"id": 2, "key": "value-2"},
//   }
//   response.With.List(v)
//
// [Return]:
//   return {
//       "data": {
//           "body": [
//               {
//               	"id": 1,
//               	"key": "value-1",
//               },
//               {
//               	"id": 2,
//               	"key": "value-2",
//               }
//           ],
//           "paging": {
//               "first": 1,
//               "before":1,
//               "current":1,
//               "next":1,
//               "last":1,
//               "limit":10,
//               "total_pages":1,
//               "total_items":2
//           }
//       },
//       "data_type": "PAGING",
//       "errno": 0,
//       "error": "",
//   }
func (o *WithManager) Paging(v interface{}, total int64, limit, page int) *Result {
	r := NewResult(TypePaging)
	r.Data = map[string]interface{}{ResultFieldForBody: v, ResultNameForPaging: NewPaginator(total, limit, page)}
	return r
}

// Success
// 返回成功数据.
//
// [Code]:
//   response.With.Success()
//
// [Return]
//   return {
//       "data": {},
//       "data_type": "OBJECT",
//       "errno": 0,
//       "error": ""
//   }
func (o *WithManager) Success() *Result {
	return o.Data(defaultOnData)
}

// 构建实例.
func (o *WithManager) init() *WithManager {
	return o
}
