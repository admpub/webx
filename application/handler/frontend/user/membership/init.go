package membership

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webfront/initialize/frontend"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 会员套餐
		agentG := u.Group(`/membership`)
		agentG.Route(`GET`, `/index`, Index).SetName(`user.membership`)
		agentG.Route(`GET,POST`, `/buy/:packageId`, Buy).SetName(`user.membership.buy`)
	})

}
