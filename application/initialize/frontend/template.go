package frontend

import (
	_ "github.com/admpub/webx/application/library/formbuilder"
	"github.com/admpub/webx/application/library/xtemplate"
	"github.com/coscms/webcore/library/ntemplate"
)

// TmplPathFixers 前台模板文件路径修正器
var TmplPathFixers = xtemplate.New(xtemplate.KindFrontend, ntemplate.NewPathAliases(), true)
