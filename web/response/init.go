// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

package response

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Code = (&CodeManager{}).init()
		With = (&WithManager{}).init()
	})
}
