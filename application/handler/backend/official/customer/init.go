package customer

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/backend/official/agent"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/complaint"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/group"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/invitation"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/level"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/prepaidcard"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/role"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/wallet"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户管理
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)
		g.Route(`GET,POST`, `/kick`, Kick)
	})
}
