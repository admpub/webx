//go:build !bindata
// +build !bindata

package page

import (
	"net/http"

	"github.com/coscms/webcore/library/httpserver"
)

func initTemplateDiskOtherFS() {
	for _, tmplDir := range httpserver.Frontend.Template.TmplDirs() { // /***/template/frontend
		templateDiskFS.Register(http.Dir(tmplDir))
	}
}
