// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

package response

import (
	"encoding/json"
)

type (
	Paginator struct {
		First      int   `json:"first" label:"First page"`
		Before     int   `json:"before" label:"Previous page"`
		Current    int   `json:"current" label:"Current page"`
		Next       int   `json:"next" label:"Next page"`
		Last       int   `json:"last" label:"Last page"`
		Limit      int   `json:"limit" label:"Count per page"`
		TotalPages int   `json:"total_pages" label:"Total pages"`
		TotalItems int64 `json:"total_items" label:"Total items"`
	}
)

const (
	defaultPaginatorLimit = 10
	defaultPaginatorPage  = 1
)

func NewPaginator(total int64, limit, page int) *Paginator {
	o := &Paginator{
		First: 1, Before: 1, Current: page, Next: 1, Last: 1,
		Limit: limit, TotalPages: 1, TotalItems: total,
	}

	if o.Limit == 0 {
		o.Limit = defaultPaginatorLimit
	}

	if o.Current == 0 {
		o.Current = defaultPaginatorPage
	}

	if t := int(total); t > limit {
		o.TotalPages = t / limit
		if t%limit > 0 {
			o.TotalPages += 1
		}
	}

	if o.Current > o.TotalPages {
		o.Current = o.TotalPages
	}

	if o.Current > 1 {
		o.Before = o.Current - 1
	}

	if o.Current < o.TotalPages {
		o.Next = o.Current + 1
	} else {
		o.Next = o.TotalPages
	}

	o.Last = o.TotalPages
	return o
}

func (o *Paginator) Json() string {
	buf, _ := json.Marshal(o)
	return string(buf)
}
