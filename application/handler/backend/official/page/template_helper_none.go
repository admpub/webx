//go:build !bindata
// +build !bindata

package page

import (
	"net/http"

	xbindata "github.com/admpub/webx/application/library/bindata"
)

func initTemplateBindata() {
	for _, tmplDir := range xbindata.FrontendTemplateDirs.TmplDirs() {
		templateDiskFS.Register(http.Dir(tmplDir))
	}
}
