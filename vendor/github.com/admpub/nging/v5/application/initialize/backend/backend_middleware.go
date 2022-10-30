package backend

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/webx-top/echo/subdomains"
)

func addGlobalFuncMap(fm map[string]interface{}) map[string]interface{} {
	fm[`AssetsURL`] = getAssetsURL
	fm[`BackendURL`] = getBackendURL
	fm[`FrontendURL`] = getFrontendURL
	fm[`Project`] = navigate.ProjectGet
	fm[`ProjectSearchIdent`] = navigate.ProjectSearchIdent
	fm[`Projects`] = navigate.ProjectListAll
	return fm
}

func getAssetsURL(paths ...string) (r string) {
	r = AssetsURLPath
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
	return r
	//return subdomains.Default.URL(r, `backend`)
}

func getFrontendURL(paths ...string) (r string) {
	r = handler.FrontendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return subdomains.Default.URL(r, `frontend`)
}
