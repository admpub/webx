package frontend

import (
	"github.com/admpub/webx/application/library/xtemplate"
	"github.com/webx-top/echo"
)

func FrontendURLFuncMW() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			FrontendURLFunc(c)
			return h.Handle(c)
		})
	}
}

func FrontendURLFunc(c echo.Context) error {
	c.SetFunc(`ThemeInfo`, func() *xtemplate.ThemeInfo {
		return TmplPathFixers.ThemeInfo(c)
	})
	return nil
}
