package customer

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	_ "github.com/admpub/webx/application/handler/backend/official/customer/complaint"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/group"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/group_package"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/invitation"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/level"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/prepaidcard"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/role"
	_ "github.com/admpub/webx/application/handler/backend/official/customer/wallet"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/registry/route"

	"github.com/admpub/nging/v5/application/handler/manager"
)

func init() {
	manager.AddUserLink(func(ctx echo.Context, user *dbschema.NgingUser) string {
		return `<a href="` + backend.URLFor(`/user/message/send`) + `?recipientType=user&recipientId=` + param.AsString(user.Id) + `" target="_blank" title="` + ctx.T(`发站内信`) + `"><i class="fa fa-regular fa-envelope"></i></a>`
	})
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
