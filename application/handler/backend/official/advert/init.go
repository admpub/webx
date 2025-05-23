package advert

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/advert`)
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)
		g.Route(`GET,POST`, `/position_index`, PositionIndex)
		g.Route(`GET,POST`, `/position_add`, PositionAdd)
		g.Route(`GET,POST`, `/position_edit`, PositionEdit)
		g.Route(`GET,POST`, `/position_delete`, PositionDelete)
	})
}
