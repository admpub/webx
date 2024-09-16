package customer

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/backend/official/customer/complaint"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/group"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/group_package"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/invitation"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/level"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/prepaidcard"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/role"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/wallet"
	"github.com/coscms/webcore/registry/route"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 客户管理
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)
		g.Route(`GET,POST`, `/kick`, Kick)
		g.Route(`GET,POST`, `/recount_file`, RecountFile)
	})
}
