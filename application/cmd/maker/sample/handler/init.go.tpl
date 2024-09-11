package {{.PkgName}}

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	{{.MakeInit}}
}
