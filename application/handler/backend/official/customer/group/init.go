package group

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户组
		g.Route(`GET,POST`, `/group/index`, Index)
		g.Route(`GET,POST`, `/group/add`, Add)
		g.Route(`GET,POST`, `/group/edit`, Edit)
		g.Route(`GET,POST`, `/group/delete`, Delete)
	})
}
