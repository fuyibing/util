// author: wsfuyibing <websearch@163.com>
// date: 2020-01-01

package util

import "sync"

func init() {
	new(sync.Once).Do(func() {
	})
}
