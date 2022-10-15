package bindata

import (
	"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/ntemplate"
)

var (
	BackendTemplateDirs  = bindata.PathAliases        //{prefix:templateDir}
	FrontendTemplateDirs = ntemplate.NewPathAliases() //{prefix:templateDir}
)
