// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

package response

var (
	With *WithManager

	defaultOnData = make(map[string]interface{})
)

type (
	WithManager struct{}
)

// Data
// return basic result.
//
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
// return error result with system code.
//
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
// return error result with specified code.
//
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
// return list result.
//
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
// return paginator results.
//
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
// return succeed result with default values.
//
//   return {
//       "data": {},
//       "data_type": "OBJECT",
//       "errno": 0,
//       "error": ""
//   }
func (o *WithManager) Success() *Result {
	return o.Data(defaultOnData)
}

func (o *WithManager) init() *WithManager {
	return o
}
