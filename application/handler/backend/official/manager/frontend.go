package manager

import (
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/initialize/frontend"
)

func FrontendReboot(ctx echo.Context) error {
	frontend.IRegister().Echo().Reset()
	echo.Fire(`webx.frontend.close`)
	frontend.InitWebServer()
	frontend.IRegister().Echo().Commit()
	return ctx.String(ctx.T(`已经重启完毕`))
}
