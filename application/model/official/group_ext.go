package official

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

var (
	//GroupTypes 组类型
	GroupTypes = echo.NewKVData().
		Add(`customer`, `客户组`).
		Add(`openapp`, `开放平台应用`).
		Add(`api`, `外部API`)
)

type GroupAndType struct {
	*dbschema.OfficialCommonGroup
	Type *echo.KV
}
