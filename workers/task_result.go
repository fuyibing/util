// author: wsfuyibing <websearch@163.com>
// date: 2022-11-16

package workers

import (
	"sync/atomic"
	"time"
)

var (
	taskResultId uint64
)

// TaskResult
// 任务结果.
type TaskResult struct {
	id uint64

	Created, Begin, Finish time.Time
	Duration, Delay        int64
	Errors                 []error
	Returned               interface{}
}

func NewTaskResult() TaskResult {
	return TaskResult{
		id:      atomic.AddUint64(&taskResultId, 1),
		Created: time.Now(),
		Errors:  make([]error, 0),
	}
}

func (o TaskResult) HasError() bool {
	return len(o.Errors) > 0
}

func (o TaskResult) Id() uint64 {
	return o.id
}
