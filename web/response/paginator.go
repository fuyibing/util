// author: wsfuyibing <websearch@163.com>
// date: 2022-10-23

package response

import (
	"encoding/json"
)

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

const (
	defaultPaginatorLimit = 10
	defaultPaginatorPage  = 1
)

// NewPaginator
// 返回分页统计.
func NewPaginator(total int64, limit, page int) *Paginator {
	// 1. 基础统计.
	o := &Paginator{
		First: 1, Before: 1, Current: page, Next: 1, Last: 1,
		Limit: limit, TotalPages: 1, TotalItems: total,
	}
	// 1.1 每页数量.
	if o.Limit == 0 {
		o.Limit = defaultPaginatorLimit
	}
	// 1.2 当前页码.
	if o.Current == 0 {
		o.Current = defaultPaginatorPage
	}

	// 2. 最大页码.
	if t := int(total); t > limit {
		o.TotalPages = t / limit
		if t%limit > 0 {
			o.TotalPages += 1
		}
	}

	// 3. 当前页.
	if o.Current > o.TotalPages {
		o.Current = o.TotalPages
	}

	// 4. 上一页.
	if o.Current > 1 {
		o.Before = o.Current - 1
	}

	// 5. 下一页.
	if o.Current < o.TotalPages {
		o.Next = o.Current + 1
	} else {
		o.Next = o.TotalPages
	}

	// 6. 末页.
	o.Last = o.TotalPages

	return o
}

// Json
// 转成JSON字符串.
func (o *Paginator) Json() string {
	buf, _ := json.Marshal(o)
	return string(buf)
}
