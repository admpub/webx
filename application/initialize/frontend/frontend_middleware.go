package frontend

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/initialize/backend"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/webx/application/library/xtemplate"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
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
	cfgCDN := config.Setting(`base`, `assetsCDN`)
	c.SetFunc(`AssetsURL`, func(paths ...string) (r string) {
		r = backend.AssetsURLPath
		if assetsCDN := cfgCDN.String(`backend`); len(assetsCDN) > 0 {
			r = assetsCDN
		}
		for _, ppath := range paths {
			r += ppath
		}
		return r
	})
	c.SetFunc(`AssetsXURL`, func(paths ...string) (r string) {
		r = AssetsURLPath
		if assetsCDN := cfgCDN.String(`frontend`); len(assetsCDN) > 0 {
			r = assetsCDN
		}
		for _, ppath := range paths {
			r += ppath
		}
		return r
	})
	c.SetFunc(`BackendURL`, func(paths ...string) (r string) {
		r = handler.BackendPrefix
		for _, ppath := range paths {
			r += ppath
		}
		return subdomains.Default.URL(r, `backend`)
	})
	c.SetFunc(`FrontendURL`, func(paths ...string) (r string) {
		r = handler.FrontendPrefix
		for _, ppath := range paths {
			r += ppath
		}
		return subdomains.Default.URL(r, `frontend`)
	})
	c.SetFunc(`SubDomain`, subdomains.Default.Get(`frontend`).TypeHost)
	c.SetFunc(`ThemeInfo`, func() *xtemplate.ThemeInfo {
		return TmplPathFixers.ThemeInfo(c)
	})
	return nil
}
