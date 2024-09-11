package application

import (
	_ "github.com/admpub/webx/application/cmd"
	_ "github.com/admpub/webx/application/handler"
	_ "github.com/admpub/webx/application/initialize"
	_ "github.com/admpub/webx/application/listener"
	_ "github.com/admpub/webx/application/registry"
	"github.com/coscms/webcore/cmd/bootconfig"
)

func init() {
	_ = bootconfig.SupportManager
	//bootconfig.SupportManager = false
	//echo.Set(`BackendPrefix`, `/admin`)
}
