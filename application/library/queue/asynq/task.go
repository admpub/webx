package asynq

import (
	"github.com/hibiken/asynq"
)

type Task = asynq.Task

func NewTask(typename string, payload []byte) *Task {
	return asynq.NewTask(typename, payload)
}
