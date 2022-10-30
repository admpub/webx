package wallet

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 资产钱包
		g.Route(`GET,POST`, `/wallet/index`, Index)
		g.Route(`GET,POST`, `/wallet/edit`, Edit)
		g.Route(`GET,POST`, `/wallet/flow`, FlowIndex)
	})
}
