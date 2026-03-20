package group_package

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 线下转账
		g.Route(`GET,POST`, `/offline_pay/index`, Index)
		g.Route(`GET,POST`, `/offline_pay/add`, Add)
		g.Route(`GET,POST`, `/offline_pay/edit`, Edit)
		g.Route(`GET,POST`, `/offline_pay/delete`, Delete)
	})
}
