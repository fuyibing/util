// author: wsfuyibing <websearch@163.com>
// date: 2022-09-30

package util

import "sync"

func init() {
    new(sync.Once).Do(func() {
        catchablePool = sync.Pool{New: func() interface{} { return (&catchable{}).init() }}
    })
}
