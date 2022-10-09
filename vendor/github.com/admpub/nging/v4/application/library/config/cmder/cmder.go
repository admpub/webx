package cmder

import (
	"io"
)

type Cmder interface {
	Init() error
	StopHistory(...string) error
	Start(writer ...io.Writer) error
	Stop() error
	Reload() error
	Restart(writer ...io.Writer) error
}

type RestartBy interface {
	RestartBy(id string, writer ...io.Writer) error
}

type StopBy interface {
	StopBy(id string) error
}

type StartBy interface {
	StartBy(id string, writer ...io.Writer) error
}
