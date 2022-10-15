//go:build bindata
// +build bindata

package page

import (
	"net/http"

	"github.com/admpub/webx/application/initialize/frontend"
	xbindata "github.com/admpub/webx/application/library/bindata"
	"github.com/admpub/webx/application/library/xtemplate"
)

func initTemplateBindata() {
	for _, tmplDir := range xbindata.FrontendTemplateDirs.TmplDirs() {
		templateDiskFS.Register(http.Dir(tmplDir))
	}
}
