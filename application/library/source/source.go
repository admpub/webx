package source

import (
	"github.com/webx-top/echo"
)

func New() *Source {
	return &Source{KVData: echo.NewKVData()}
}

type Source struct {
	*echo.KVData
}
