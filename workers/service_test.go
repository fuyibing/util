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

func TestNewService(t *testing.T) {

	var (
		ctx, cancel = context.WithCancel(context.TODO())
		err         error
		s           = NewService()
	)

	go func() {
		time.Sleep(time.Second)

		for i := 0; i < 3; i++ {
			func(index int) {
				c := trace.New()
				log.Infofc(c, "task: %d", index)

				if e := s.Add(NewTask().SetContext(c).SetFinish(func(ctx context.Context, res TaskResult) {
					log.Infofc(ctx, "task finish")
				}).SetHandler(func(ctx context.Context) interface{} {
					log.Infofc(ctx, "task handle")
					time.Sleep(time.Second)
					return nil
				})); e != nil {
					log.Errorfc(c, "task add: %v", e)
				}
			}(i)
		}

		cancel()
	}()

	t.Logf("testing.Begin")

	if err = s.SetParallel(1).Start(ctx); err != nil {
		t.Errorf("start error: %v", err)
	}

	t.Logf("testing.End")
}
