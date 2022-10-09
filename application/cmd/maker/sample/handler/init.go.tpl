package {{.PkgName}}

import (
	"github.com/admpub/nging/v3/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	{{.MakeInit}}
}
