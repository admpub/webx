package application

import (
	_ "github.com/admpub/nging/v4/application"
	"github.com/admpub/nging/v4/application/cmd/bootconfig"
	_ "github.com/admpub/nging/v4/application/library/sqlite"
	_ "github.com/admpub/webx/application/cmd"
	_ "github.com/admpub/webx/application/handler"
	_ "github.com/admpub/webx/application/initialize"
	_ "github.com/admpub/webx/application/listener"
	_ "github.com/admpub/webx/application/registry"
)

func init() {
	_ = bootconfig.SupportManager
	//bootconfig.SupportManager = false
	//echo.Set(`BackendPrefix`, `/admin`)
}
