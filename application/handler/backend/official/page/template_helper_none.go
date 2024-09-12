//go:build !bindata
// +build !bindata

package page

import (
	"net/http"

	"github.com/coscms/webfront/initialize/frontend"
)

func initTemplateDiskOtherFS() {
	for _, tmplDir := range frontend.TmplPathFixers.TmplDirs() { // /***/template/frontend
		templateDiskFS.Register(http.Dir(tmplDir))
	}
}
