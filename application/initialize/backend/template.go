package backend

import (
	_ "github.com/admpub/webx/application/library/setup"
	"github.com/admpub/webx/application/library/xtemplate"
)

// TmplPathFixers 模版路径 {subdir:func}
var TmplPathFixers = xtemplate.New(`backend`)

func init() {
	xtemplate.Register(`backend`, TmplPathFixers)
}
