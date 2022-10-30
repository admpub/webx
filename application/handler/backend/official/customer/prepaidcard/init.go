package prepaidcard

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
)

func init() {
	handler.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/customer`)
		// 充值卡
		g.Route(`GET,POST`, `/prepaid_card/index`, Index)
		g.Route(`GET,POST`, `/prepaid_card/add`, Add)
		g.Route(`GET,POST`, `/prepaid_card/edit`, Edit)
		g.Route(`GET,POST`, `/prepaid_card/delete`, Delete)
	})
}
