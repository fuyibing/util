// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package request

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Validate = (&Validator{}).init()
	})
}
