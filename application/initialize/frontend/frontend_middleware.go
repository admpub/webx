package frontend

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/initialize/backend"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/webx/application/library/xtemplate"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
)

func addGlobalFuncMap(fm map[string]interface{}) map[string]interface{} {
	fm[`AssetsURL`] = getAssetsURL
	fm[`BackendURL`] = getBackendURL
	fm[`FrontendURL`] = getFrontendURL
	fm[`SubDomain`] = subdomains.Default.Get(`frontend`).TypeHost
	return fm
}

func getAssetsURL(paths ...string) (r string) {
	r = backend.AssetsURLPath
	if assetsCDN := config.Setting(`base`, `assetsCDN`).String(`backend`); len(assetsCDN) > 0 {
		r = assetsCDN
	}
	for _, ppath := range paths {
		r += ppath
	}
	return r
}

func getAssetsXURL(paths ...string) (r string) {
	r = AssetsURLPath
	if assetsCDN := config.Setting(`base`, `assetsCDN`).String(`frontend`); len(assetsCDN) > 0 {
		r = assetsCDN
	}
	for _, ppath := range paths {
		r += ppath
	}
	return r
}

func getBackendURL(paths ...string) (r string) {
	r = handler.BackendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return subdomains.Default.URL(r, `backend`)
}

func getFrontendURL(paths ...string) (r string) {
	r = handler.FrontendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return subdomains.Default.URL(r, `frontend`)
}

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
