package article

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/article`)
		// 文章
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)

		// 分类
		g.Route(`GET,POST`, `/category`, CategoryIndex)
		g.Route(`GET,POST`, `/category_add`, CategoryAdd)
		g.Route(`GET,POST`, `/category_edit`, CategoryEdit)
		g.Route(`GET,POST`, `/category_delete`, CategoryDelete)

		// 友情链接
		g.Route(`GET,POST`, `/friendlink_index`, FriendlinkIndex)
		g.Route(`GET,POST`, `/friendlink_add`, FriendlinkAdd)
		g.Route(`GET,POST`, `/friendlink_edit`, FriendlinkEdit)
		g.Route(`GET,POST`, `/friendlink_delete`, FriendlinkDelete)
	})
}
