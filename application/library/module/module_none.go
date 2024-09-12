//go:build !bindata
// +build !bindata

package module

import (
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/bindata"
	"github.com/coscms/webcore/library/module"
)

func SetFrontendTemplate(key string, templatePath string) {
	module.SetTemplate(frontend.TmplPathFixers.PathAliases, key, templatePath)
}

func SetFrontendAssets(assetsPath string) {
	module.SetAssets(bindata.StaticOptions, assetsPath)
}
