// author: wsfuyibing <websearch@163.com>
// date: 2022-11-17

package workers

import (
	"context"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"testing"
	"time"
)

func TestNewBatch(t *testing.T) {
	c := trace.New()
	log.Infofc(c, "TestNewBatch")

	b := NewBatch().SetParallel(3)

	for i := 0; i < 100; i++ {
		b.Add(NewTask().SetContext(c).SetHandler(func(ctx context.Context) interface{} {
			time.Sleep(time.Millisecond * 10)
			return nil
		}).SetFinish(func(ctx context.Context, res TaskResult) {
			log.Infofc(ctx, "task finish: id=%d, error=%v", res.Id(), res.HasError())
		}))
	}

	total, success := b.Run()
	log.Infofc(c, "batch: total=%d, success=%d", total, success)
}
