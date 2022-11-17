// author: wsfuyibing <websearch@163.com>
// date: 2022-11-16

package workers

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		batchPool = sync.Pool{New: func() interface{} { return (&batch{}).init() }}
		servicePool = sync.Pool{New: func() interface{} { return (&service{}).init() }}
		taskPool = sync.Pool{New: func() interface{} { return (&task{}).init() }}
	})
}
