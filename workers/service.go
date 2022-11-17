// author: wsfuyibing <websearch@163.com>
// date: 2022-11-17

package workers

import (
	"context"
	"fmt"
	"github.com/fuyibing/log/v3"
	"sync"
	"sync/atomic"
	"time"
)

var (
	serviceId   uint64
	servicePool sync.Pool
)

const (
	serviceParallel int32 = 10
)

type (
	Service interface {
		Add(task Task) (err error)
		Release()
		SetParallel(parallel int32) Service
		Start(ctx context.Context) (err error)
	}

	service struct {
		acquires, id uint64
		mu           sync.RWMutex

		cancel                context.CancelFunc
		ctx                   context.Context
		parallel, concurrency int32
		started               bool
		taskChan              chan Task
		taskIndex             uint64
		taskMapper            map[uint64]Task
		total, success        int64
	}
)

// NewService
// 获取服务实例.
func NewService() Service {
	return servicePool.Get().(*service).before()
}

// Add
// 添加任务.
func (o *service) Add(task Task) (err error) {
	// 捕获异常.
	// 当 Add 方法在收到退出信号之后触发.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// 健康检查.
	if func() bool {
		o.mu.RLock()
		defer o.mu.RUnlock()
		if o.started && o.ctx != nil && o.ctx.Err() == nil {
			return false
		}
		return true
	}() {
		err = fmt.Errorf("service not started or stopping")
		return
	}

	// 中间状态.
	if o.taskChan == nil {
		time.Sleep(time.Millisecond)
		return o.Add(task)
	}

	// 发送消息.
	o.taskChan <- task
	return
}

// Release
// 释放实例.
func (o *service) Release() {
	o.after()
	servicePool.Put(o)
}

// SetParallel
// 设置最大并行任务数.
func (o *service) SetParallel(parallel int32) Service {
	if n := atomic.LoadInt32(&o.parallel); n != parallel {
		atomic.StoreInt32(&o.parallel, parallel)
		log.Infof("[worker][service] change parallel from %d to %d.", n, parallel)

		// 增加并行.
		if parallel > n && func() bool {
			o.mu.RLock()
			defer o.mu.RUnlock()
			return o.ctx != nil && o.ctx.Err() == nil
		}() {
			for i := 0; i < int(parallel-n); i++ {
				go o.pop()
			}
		}
	}

	return o
}

// Start
// 启动服务.
func (o *service) Start(ctx context.Context) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	o.mu.Lock()

	// 重复启动.
	if o.started {
		o.mu.Unlock()
		err = fmt.Errorf("service started already")
		return
	}

	// 准备启动.
	o.ctx, o.cancel = context.WithCancel(ctx)
	o.started = true
	o.taskChan = make(chan Task)
	o.mu.Unlock()
	log.Infof("[worker][service] start")

	// 退出服务.
	defer func() {
		// 服务异常.
		if r := recover(); r != nil {
			log.Panicf("service panic: %v", r)
		}

		// 强制取消.
		if o.ctx.Err() == nil {
			o.cancel()
		}

		// 关闭通道.
		close(o.taskChan)
		o.mu.Lock()
		o.taskChan = nil
		o.mu.Unlock()

		// 等待完成.
		o.wait()

		// 恢复字段.
		log.Infof("[worker][service] stopped")
		o.mu.Lock()
		o.ctx = nil
		o.cancel = nil
		o.mu.Unlock()
	}()

	// 待等消息.
	for {
		select {
		case x := <-o.taskChan:
			go o.push(x)
		case <-o.ctx.Done():
			return
		}
	}
}

// /////////////////////////////////////////////////////////////
// Pool access operations
// /////////////////////////////////////////////////////////////

func (o *service) pop() {
	var (
		concurrency = atomic.AddInt32(&o.concurrency, 1)
		tasker      Task
	)

	// 任务限流.
	if concurrency > o.parallel {
		atomic.AddInt32(&o.concurrency, -1)
		return
	}

	// 处理结束.
	defer func() {
		atomic.AddInt32(&o.concurrency, -1)

		// 继续取出.
		if func() bool {
			o.mu.RLock()
			defer o.mu.RUnlock()
			return len(o.taskMapper) > 0
		}() {
			o.pop()
		}
	}()

	// 取出数据.
	tasker = func() Task {
		o.mu.Lock()
		defer o.mu.Unlock()
		for k, v := range o.taskMapper {
			delete(o.taskMapper, k)
			return v
		}
		return nil
	}()

	// 退出处理.
	// 从内存中未取到消息(即消息取完了).
	if tasker == nil {
		return
	}

	// 处理任务.
	atomic.AddInt64(&o.total, 1)
	if tasker.Run() {
		atomic.AddInt64(&o.success, 1)
	}
}

func (o *service) push(task Task) {
	i := atomic.AddUint64(&o.taskIndex, 1)

	// 写入缓冲.
	o.mu.Lock()
	o.taskMapper[i] = task
	o.mu.Unlock()

	// 取出任务.
	o.pop()
}

func (o *service) wait() {
	o.mu.RLock()

	if len(o.taskMapper) > 0 || atomic.LoadInt32(&o.concurrency) > 0 {
		o.mu.RUnlock()

		log.Debugf("[worker][service] waiting tasks finish")
		time.Sleep(time.Millisecond * 100)
		o.wait()
		return
	}

	o.mu.RUnlock()
	log.Infof("[worker][service] task finished: total=%d, success=%d", o.total, o.success)
}

// /////////////////////////////////////////////////////////////
// Pool instance operations
// /////////////////////////////////////////////////////////////

func (o *service) after() {
	o.taskMapper = nil
}

func (o *service) before() *service {
	atomic.AddUint64(&o.acquires, 1)

	o.concurrency = 0
	o.parallel = serviceParallel
	o.started = false
	o.taskIndex = 0
	o.taskMapper = make(map[uint64]Task)
	o.total = 0
	o.success = 0
	return o
}

func (o *service) init() *service {
	o.id = atomic.AddUint64(&serviceId, 1)
	o.mu = sync.RWMutex{}
	return o
}
