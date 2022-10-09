package manager

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/registry/navigate"
)

func init() {
	handler.RegisterToGroup(`/manager`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/message/index`, MessageIndex)
		g.Route(`GET,POST`, `/message/delete`, MessageDelete)
		g.Route(`GET,POST`, `/message/view/:id`, MessageView)

		// 导航菜单
		g.Route(`GET,POST`, `/navigate/index`, NavigateIndex)
		g.Route(`GET,POST`, `/navigate/add`, NavigateAdd)
		g.Route(`GET,POST`, `/navigate/edit`, NavigateEdit)
		g.Route(`GET,POST`, `/navigate/delete`, NavigateDelete)

		g.Route(`GET,POST`, `/frontend/reboot`, FrontendReboot)
		g.Route(`GET,POST`, `/frontend/route_page`, FrontendRoutePage)
		g.Route(`GET,POST`, `/frontend/route_page_add`, FrontendRoutePageAdd)
		g.Route(`GET,POST`, `/frontend/route_page_edit`, FrontendRoutePageEdit)
		g.Route(`GET,POST`, `/frontend/route_page_delete`, FrontendRoutePageDelete)
	})

	(*navigate.TopNavigate)[0].Children.Add(
		-1,
		TopNavigate...,
	)
}
