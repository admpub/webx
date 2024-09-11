package level

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户等级
		g.Route(`GET,POST`, `/level/index`, Index)
		g.Route(`GET,POST`, `/level/add`, Add)
		g.Route(`GET,POST`, `/level/edit`, Edit)
		g.Route(`GET,POST`, `/level/delete`, Delete)
	})
}
