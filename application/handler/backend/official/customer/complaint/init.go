package complaint

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户组
		g.Route(`GET,POST`, `/complaint/index`, Index)
		g.Route(`GET,POST`, `/complaint/edit`, Edit)
		g.Route(`GET,POST`, `/complaint/delete`, Delete)
	})
}
