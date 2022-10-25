// author: wsfuyibing <websearch@163.com>
// date: 2022-10-23

package response

type (

	// Paginator
	// 分页.
	Paginator struct {
		First      int   `json:"first"`
		Before     int   `json:"before"`
		Current    int   `json:"current"`
		Next       int   `json:"next"`
		Last       int   `json:"last"`
		Limit      int   `json:"limit"`
		TotalPages int   `json:"total_pages"`
		TotalItems int64 `json:"total_items"`
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
