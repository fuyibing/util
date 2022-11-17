// author: wsfuyibing <websearch@163.com>
// date: 2022-11-17

package workers

import (
	"context"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"testing"
)

func TestNewTask(t *testing.T) {
	c := trace.New()
	log.Infofc(c, "TestNewTask")

	for i := 0; i < 2; i++ {
		NewTask().SetContext(c).SetHandler(func(ctx context.Context) interface{} {
			log.Infofc(ctx, "task handler")
			return "task handle finish"
		}).SetFinish(func(ctx context.Context, res TaskResult) {
			log.Infofc(ctx, "task finish: error=%v, %+v", res.HasError(), res)
		}).Run()
	}
}
