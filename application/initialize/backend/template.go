package backend

import (
	_ "github.com/admpub/webx/application/library/setup"
	"github.com/admpub/webx/application/library/xtemplate"
	"github.com/coscms/webcore/library/bindata"
)

// TmplPathFixers 后台模板文件路径修正器
var TmplPathFixers = xtemplate.New(`backend`, bindata.PathAliases)

func init() {
	xtemplate.Register(`backend`, TmplPathFixers)
}
