package machinery

import "github.com/RichardKnop/machinery/v1/tasks"

type (
	Task     = tasks.Signature
	TaskArg  = tasks.Arg
	TaskArgs []TaskArg
)

func Arg(dataType string, dataValue interface{}, name ...string) TaskArg {
	a := TaskArg{Type: dataType, Value: dataValue}
	if len(name) > 0 {
		a.Name = name[0]
	}
	return a
}

func Args() TaskArgs {
	return TaskArgs{}
}

func (t *TaskArgs) Add(dataType string, dataValue interface{}, name ...string) *TaskArgs {
	*t = append(*t, Arg(dataType, dataValue, name...))
	return t
}

func (t *TaskArgs) AddArg(args ...TaskArg) *TaskArgs {
	*t = append(*t, args...)
	return t
}

func (t *TaskArgs) Slice() TaskArgs {
	return *t
}

// NewTask creates a new task signature
func NewTask(name string, args ...TaskArg) *Task {
	t, _ := tasks.NewSignature(name, args)
	return t
}
