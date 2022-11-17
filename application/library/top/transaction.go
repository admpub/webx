package top

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func Tx(ctx echo.Context) factory.Transactioner {
	return ctx.Transaction().(*factory.Param).Trans()
}

func Dump(m interface{}) {
	echo.Dump(m)
}
