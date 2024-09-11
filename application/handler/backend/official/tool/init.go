package tool

import (
	"github.com/webx-top/echo"

	_ "github.com/coscms/webcore"
	"github.com/coscms/webcore/registry/navigate"
	"github.com/coscms/webcore/registry/route"
)

func init() {
	route.RegisterToGroup(`/tool`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/area/index`, AreaIndex)
		g.Route(`GET,POST`, `/area/edit`, AreaEdit)
		g.Route(`GET,POST`, `/area/add`, AreaAdd)
		g.Route(`GET,POST`, `/area/delete`, AreaDelete)
		g.Route(`GET,POST`, `/area/group_index`, AreaGroupIndex)
		g.Route(`GET,POST`, `/area/group_edit`, AreaGroupEdit)
		g.Route(`GET,POST`, `/area/group_add`, AreaGroupAdd)
		g.Route(`GET,POST`, `/area/group_delete`, AreaGroupDelete)

		sensitive := g.Group(`/sensitive`)
		// 敏感词管理
		sensitive.Route(`GET,POST`, `/index`, sensitiveIndex)
		sensitive.Route(`GET,POST`, `/add`, sensitiveAdd)
		sensitive.Route(`GET,POST`, `/edit`, sensitiveEdit)
		sensitive.Route(`GET,POST`, `/delete`, sensitiveDelete)

		// 中文分词
		g.Route(`GET,POST`, `/segment`, Segment)
	})

	(*navigate.TopNavigate)[1].Children.Add(
		-1,
		TopNavigate...,
	)
}
