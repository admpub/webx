package group_package

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 等级组套餐
		g.Route(`GET,POST`, `/group_package/index`, Index)
		g.Route(`GET,POST`, `/group_package/add`, Add)
		g.Route(`GET,POST`, `/group_package/edit`, Edit)
		g.Route(`GET,POST`, `/group_package/delete`, Delete)
	})
}
