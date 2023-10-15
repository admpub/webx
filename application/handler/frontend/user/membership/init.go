package membership

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/admpub/webx/application/initialize/frontend"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 会员套餐
		agentG := u.Group(`/membership`)
		agentG.Route(`GET`, `/index`, Index)
		agentG.Route(`GET,POST`, `/buy/:packageId`, Buy)
	})

}
