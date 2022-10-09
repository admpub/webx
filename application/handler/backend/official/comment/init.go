package comment

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/article`)
		// 评论列表
		g.Route(`GET,POST`, `/comment/index`, Index)
		g.Route(`GET,POST`, `/comment/list`, List)
		g.Route(`GET,POST`, `/comment/add`, Add)
		g.Route(`GET,POST`, `/comment/edit`, Edit)
		g.Route(`GET,POST`, `/comment/delete`, Delete)
	})
}
