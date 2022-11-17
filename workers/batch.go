// author: wsfuyibing <websearch@163.com>
// date: 2022-11-17

package workers

import (
	"sync"
	"sync/atomic"
)

var (
	batchId   uint64
	batchPool sync.Pool
)

const (
	batchParallel = 10
)

type (
	// Batch
	// 批量任务处理.
	Batch interface {
		Add(task Task) Batch
		Run() (total, success int64)
		SetParallel(parallel int) Batch
	}

	// 批量处理.
	batch struct {
		acquires, id uint64

		mu         sync.RWMutex
		parallel   int
		taskIndex  uint64
		taskMapper map[uint64]Task

		total, success int64
	}
)

// NewBatch
// 从池中取出实例.
func NewBatch() Batch {
	return batchPool.Get().(*batch).before()
}

// Add
// 添加任务.
func (o *batch) Add(task Task) Batch {
	i := atomic.AddUint64(&o.taskIndex, 1)

	o.mu.Lock()
	defer o.mu.Unlock()

	o.taskMapper[i] = task
	return o
}

// Run
// 批量执行.
func (o *batch) Run() (total, success int64) {
	// 计算任务
	n := func() int {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return len(o.taskMapper)
	}()

	// 最大并行.
	if n > o.parallel {
		n = o.parallel
	}

	// 并行处理.
	wait := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			o.handle()
		}()
	}
	wait.Wait()

	// 处理完成.
	o.after()
	batchPool.Put(o)
	return o.total, o.success
}

// SetParallel
// 设置最大并行.
func (o *batch) SetParallel(parallel int) Batch {
	o.parallel = parallel
	return o
}

// /////////////////////////////////////////////////////////////
// Pool instance operations
// /////////////////////////////////////////////////////////////

func (o *batch) handle() {
	// 取出任务.
	x := func() Task {
		o.mu.Lock()
		defer o.mu.Unlock()
		for k, v := range o.taskMapper {
			delete(o.taskMapper, k)
			return v
		}
		return nil
	}()

	// 任务取完.
	if x == nil {
		return
	}

	// 处理过程.
	atomic.AddInt64(&o.total, 1)
	if x.Run() {
		atomic.AddInt64(&o.success, 1)
	}

	// 继续处理.
	// 直到全部任务取出完成.
	o.handle()
}

// /////////////////////////////////////////////////////////////
// Pool instance operations
// /////////////////////////////////////////////////////////////

func (o *batch) after() {
	o.taskMapper = nil
}

func (o *batch) before() *batch {
	atomic.AddUint64(&o.acquires, 1)

	o.parallel = batchParallel
	o.taskIndex = 0
	o.taskMapper = make(map[uint64]Task, 0)
	return o
}

func (o *batch) init() *batch {
	o.id = atomic.AddUint64(&batchId, 1)
	o.mu = sync.RWMutex{}
	return o
}
