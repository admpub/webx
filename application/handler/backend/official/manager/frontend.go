package manager

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webfront/initialize/frontend"
)

func FrontendReboot(ctx echo.Context) error {
	frontend.IRegister().Echo().Reset()
	echo.Fire(`webx.frontend.close`)
	frontend.InitWebServer()
	frontend.IRegister().Echo().Commit()
	echo.Fire(`webx.frontend.reboot`)
	return ctx.String(ctx.T(`已经重启完毕`))
}
