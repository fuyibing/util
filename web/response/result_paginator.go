// author: wsfuyibing <websearch@163.com>
// date: 2022-10-23

package response

type (

	// Paginator
	// 分页.
	Paginator struct {
		First      int   `json:"first" label:"首页"`
		Before     int   `json:"before" label:"上页"`
		Current    int   `json:"current" label:"本页"`
		Next       int   `json:"next" label:"下页"`
		Last       int   `json:"last" label:"末页"`
		Limit      int   `json:"limit" label:"每页数量"`
		TotalPages int   `json:"total_pages" label:"总页数"`
		TotalItems int64 `json:"total_items" label:"总条数"`
	}
)

// NewPaginator
// 返回分页统计.
func NewPaginator(total int64, limit, page int) *Paginator {
	// 1. 创建实例.
	p := &Paginator{
		First:      1,
		Before:     1,
		Current:    page,
		Next:       1,
		Last:       1,
		Limit:      limit,
		TotalPages: 1,
		TotalItems: total,
	}

	return p
}
