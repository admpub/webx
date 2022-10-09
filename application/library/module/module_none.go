//go:build !bindata
// +build !bindata

package module

import (
	"github.com/admpub/nging/v4/application/library/module"
	"github.com/admpub/webx/application/library/bindata"
)

func SetFrontendTemplate(key string, templatePath string) {
	module.SetTemplate(bindata.FrontendTemplateDirs, key, templatePath)
}

func SetFrontendAssets(assetsPath string) {
	module.SetAssets(bindata.StaticOptions, assetsPath)
}
