package tags

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		// 总标签
		g.Route(`GET,POST`, `/tags/index`, Index)
		//g.Route(`GET,POST`, `/tags/add`, Add)
		g.Route(`GET,POST`, `/tags/edit`, Edit)
		g.Route(`GET,POST`, `/tags/delete`, Delete)
	})
}
