package wallet

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 资产钱包
		g.Route(`GET,POST`, `/wallet/index`, Index)
		g.Route(`GET,POST`, `/wallet/edit`, Edit)
		g.Route(`GET,POST`, `/wallet/flow`, FlowIndex)
	})
}
