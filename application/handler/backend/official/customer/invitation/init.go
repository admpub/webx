package invitation

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 邀请码
		g.Route(`GET,POST`, `/invitation/index`, Index)
		g.Route(`GET,POST`, `/invitation/add`, Add)
		g.Route(`GET,POST`, `/invitation/edit`, Edit)
		g.Route(`GET,POST`, `/invitation/delete`, Delete)
		g.Route(`GET,POST`, `/invitation/customer_list`, CustomerList)
	})
}
