//go:build !bindata
// +build !bindata

package page

import (
	"net/http"

	xbindata "github.com/admpub/webx/application/library/bindata"
)

func initTemplateDiskOtherFS() {
	for _, tmplDir := range xbindata.FrontendTemplateDirs.TmplDirs() { // /***/template/frontend
		templateDiskFS.Register(http.Dir(tmplDir))
	}
}
