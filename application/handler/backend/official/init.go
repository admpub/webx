package official

import (
	backendRoute "github.com/admpub/nging/v5/application/registry/route"
	"github.com/admpub/webx/application/library/top"
	"github.com/webx-top/echo"
)

func init() {
	backendRoute.IRegister().Use(func(h echo.Handler) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.SetFunc(`OutputContent`, top.OutputContent)
			return h.Handle(ctx)
		}
	})
}
