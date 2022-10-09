package asynq

import (
	"context"

	"github.com/hibiken/asynq"
)

var DefaultWorkerConfig = asynq.Config{
	// Specify how many concurrent workers to use
	Concurrency: 10,
	// You can optionally create multiple queues with different priority.
	Queues: map[string]int{
		"critical": 6,
		"default":  3,
		"low":      1,
	},
	// See the godoc for other configuration options
}

func (a *Asynq) ServeMux() *asynq.ServeMux {
	if a.serverHandler == nil {
		a.serverHandler = asynq.NewServeMux()
	}
	return a.serverHandler
}

func (a *Asynq) Use(mws ...asynq.MiddlewareFunc) *Asynq {
	a.ServeMux().Use(mws...)
	return a
}

func (a *Asynq) Handle(pattern string, handler asynq.Handler) *Asynq {
	a.ServeMux().Handle(pattern, handler)
	return a
}

func (a *Asynq) HandleFunc(pattern string, handler func(context.Context, *asynq.Task) error) *Asynq {
	a.ServeMux().HandleFunc(pattern, handler)
	return a
}

func (a *Asynq) StartWorker(configs ...*asynq.Config) error {
	if len(configs) > 0 {
		a.serverConfig = configs[0]
	}
	return a.Server().Run(a.serverHandler)
}

func (a *Asynq) newServer() *asynq.Server {
	var config asynq.Config
	if a.serverConfig != nil {
		config = *a.serverConfig
	} else {
		config = DefaultWorkerConfig
	}
	return asynq.NewServer(a.redisOptions, config)
}
