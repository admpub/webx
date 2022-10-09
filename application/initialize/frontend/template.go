package frontend

import (
	_ "github.com/admpub/webx/application/library/formbuilder"
	"github.com/admpub/webx/application/library/xtemplate"
)

// TmplPathFixers 模版路径 {subdir:func}
var TmplPathFixers = xtemplate.New(`frontend`)

func init() {
	xtemplate.Register(`frontend`, TmplPathFixers)
}
