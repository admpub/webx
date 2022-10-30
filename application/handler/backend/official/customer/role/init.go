package role

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户角色管理
		g.Route(`GET,POST`, `/role/index`, Index)
		g.Route(`GET,POST`, `/role/add`, Add)
		g.Route(`GET,POST`, `/role/edit`, Edit)
		g.Route(`GET,POST`, `/role/delete`, Delete)
	})
}
