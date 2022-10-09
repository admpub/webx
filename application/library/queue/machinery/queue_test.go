package machinery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/webx-top/echo/testing/test"

	exampletasks "github.com/admpub/webx/application/library/queue/machinery/example"
)

func TestQueue(t *testing.T) {
	cfg := &config.Config{
		//Broker:`amqp://guest:guest@localhost:5672/`,
		Broker: `redis://localhost:6379`,
		//Broker: `https://sqs.us-west-2.amazonaws.com/123456789012`,
		//Broker: `amqp://guest:guest@localhost:5672/`,
		DefaultQueue:  `machinery_tasks`,
		ResultBackend: `redis://localhost:6379`,
		//ResultBackend:`memcache://localhost:11211`,
		//ResultBackend:`mongodb://localhost:27017`,
		ResultsExpireIn: 3600000,
		// Redis: &config.RedisConfig{
		// 	MaxIdle:                3,
		// 	IdleTimeout:            240,
		// 	ReadTimeout:            15,
		// 	WriteTimeout:           15,
		// 	ConnectTimeout:         15,
		// 	NormalTasksPollPeriod:  1000,
		// 	DelayedTasksPollPeriod: 500,
		// },
	}
	server, err := NewServer(map[string]interface{}{
		"add":               exampletasks.Add,
		"multiply":          exampletasks.Multiply,
		"sum_ints":          exampletasks.SumInts,
		"sum_floats":        exampletasks.SumFloats,
		"concat":            exampletasks.Concat,
		"split":             exampletasks.Split,
		"panic_task":        exampletasks.PanicTask,
		"long_running_task": exampletasks.LongRunningTask,
	}, cfg)
	if err != nil {
		panic(err)
	}
	go func() {
		err := server.StartWorker()
		if err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	asyncResult, err := server.Send(ctx,
		NewTask(
			`add`,
			Arg(`int64`, 1),
			Arg(`int64`, 1),
		),
	)
	if err != nil {
		panic(err)
	}
	results, err := asyncResult[0].Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		panic(fmt.Errorf("Getting task result failed with error: %s", err.Error()))
	}
	test.Eq(t, `2`, tasks.HumanReadableResults(results))
	test.Eq(t, int64(2), results[0].Interface())
}
