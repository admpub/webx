package top

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func Tx(ctx echo.Context) *factory.Transaction {
	t := ctx.Transaction().(*factory.Param).Trans()
	return t
}

func Dump(m interface{}) {
	echo.Dump(m)
}
