package frontend

import (
	"sync"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/initialize/backend"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/webx/application/library/top"
	"github.com/admpub/webx/application/middleware/sessdata"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/subdomains"
)

var (
	AssetsURLPath = `/public/assets/frontend`
	tplFuncMap    map[string]interface{}
	tplOnce       sync.Once
)

func initTplFuncMap() {
	tplFuncMap = addGlobalFuncMap(tplfunc.New())
}

func GlobalFuncMap() map[string]interface{} {
	tplOnce.Do(initTplFuncMap)
	return tplFuncMap
}

func init() {
	tplfunc.TplFuncMap[`ImageProxyURL`] = sessdata.ImageProxyURL
	tplfunc.TplFuncMap[`ResizeImageURL`] = sessdata.ResizeImageURL
	tplfunc.TplFuncMap[`AbsoluteURL`] = sessdata.AbsoluteURL
	tplfunc.TplFuncMap[`PictureHTML`] = sessdata.PictureWithDefaultHTML
	tplfunc.TplFuncMap[`OutputContent`] = sessdata.OutputContent
	tplfunc.TplFuncMap[`Config`] = config.FromDB
	tplfunc.TplFuncMap[`StarsSlice`] = top.StarsSlice
	tplfunc.TplFuncMap[`StarsSlicex`] = top.StarsSlicex
}

func addGlobalFuncMap(fm map[string]interface{}) map[string]interface{} {
	fm[`AssetsURL`] = getAssetsURL
	fm[`BackendURL`] = getBackendURL
	fm[`AssetsXURL`] = getAssetsXURL
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
