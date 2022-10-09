package machinery

import (
	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/admpub/log"
)

// Worker 消费者

type WorkerConfig struct {
	ConsumerTag     string
	Concurrency     int
	errorHandler    func(err error)
	preTaskHandler  func(signature *tasks.Signature)
	postTaskHandler func(signature *tasks.Signature)
}

type WorkerOption func(*WorkerConfig)

func ConsumerTag(tag string) WorkerOption {
	return func(config *WorkerConfig) {
		config.ConsumerTag = tag
	}
}

func Concurrency(concurrency int) WorkerOption {
	return func(config *WorkerConfig) {
		config.Concurrency = concurrency
	}
}

func ErrorHandler(errorHandler func(err error)) WorkerOption {
	return func(config *WorkerConfig) {
		config.errorHandler = errorHandler
	}
}

func PreTaskHandler(preTaskHandler func(signature *tasks.Signature)) WorkerOption {
	return func(config *WorkerConfig) {
		config.preTaskHandler = preTaskHandler
	}
}

func PostTaskHandler(postTaskHandler func(signature *tasks.Signature)) WorkerOption {
	return func(config *WorkerConfig) {
		config.postTaskHandler = postTaskHandler
	}
}

func (s *Server) StartWorker(options ...WorkerOption) error {
	c := &WorkerConfig{
		ConsumerTag: `machinery_worker`,
		Concurrency: 10,
		// Here we inject some custom code for error handling,
		// start and end of task hooks, useful for metrics for example.
		errorHandler: func(err error) {
			log.Error("I am an error handler: ", err.Error())
		},
		preTaskHandler: func(signature *tasks.Signature) {
			log.Debug("I am a start of task handler for: ", signature.Name)
		},
		postTaskHandler: func(signature *tasks.Signature) {
			log.Debug("I am an end of task handler for: ", signature.Name)
		},
	}
	for _, option := range options {
		option(c)
	}
	server := s.Server

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(c.ConsumerTag, c.Concurrency)
	worker.SetPostTaskHandler(c.postTaskHandler)
	worker.SetErrorHandler(c.errorHandler)
	worker.SetPreTaskHandler(c.preTaskHandler)

	return worker.Launch()
}
