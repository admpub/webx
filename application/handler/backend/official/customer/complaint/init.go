package complaint

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户组
		g.Route(`GET,POST`, `/complaint/index`, Index)
		g.Route(`GET,POST`, `/complaint/edit`, Edit)
		g.Route(`GET,POST`, `/complaint/delete`, Delete)
	})
}
