// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

package request

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Validate = (&Validator{}).init()
	})
}
