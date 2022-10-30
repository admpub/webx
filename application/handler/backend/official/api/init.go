package api

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/handler/backend/official/api/account"
	"github.com/webx-top/echo"
)

func init() {

	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/api`)
		// 外部接口账号
		g.Route(`GET,POST`, `/account/index`, account.Index)
		g.Route(`GET,POST`, `/account/add`, account.Add)
		g.Route(`GET,POST`, `/account/edit/:id`, account.Edit)
		g.Route(`GET,POST`, `/account/delete/:id`, account.Delete)
	})
}
