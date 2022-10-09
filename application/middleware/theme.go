package middleware

import (
	"github.com/admpub/webx/application/library/xtemplate"
	"github.com/webx-top/echo"
)

// SetTheme 设置主题
func SetTheme(theme string) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Internal().Set(`theme`, theme)
			return h.Handle(c)
		}
	}
}

// UseTheme 使用系统设置的主题
func UseTheme(fn func(echo.Context) *xtemplate.ThemeInfo) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			themeCfg := fn(c)
			theme := themeCfg.Name
			if len(theme) > 0 {
				c.Internal().Set(`theme.config`, themeCfg)
				c.Internal().Set(`theme`, theme)
			}
			return h.Handle(c)
		}
	}
}
