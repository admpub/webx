package agent

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/agent`)
		// 代理商
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)

		// 代理等级
		g.Route(`GET,POST`, `/level_index`, LevelIndex)
		g.Route(`GET,POST`, `/level_add`, LevelAdd)
		g.Route(`GET,POST`, `/level_edit`, LevelEdit)
		g.Route(`GET,POST`, `/level_delete`, LevelDelete)

		// 代理产品
		g.Route(`GET,POST`, `/product_index`, ProductIndex)
		g.Route(`GET,POST`, `/product_add`, ProductAdd)
		g.Route(`GET,POST`, `/product_edit`, ProductEdit)
		g.Route(`GET,POST`, `/product_delete`, ProductDelete)
	})
}
