package wallet

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 钱包
		g := u.Group(`/wallet`)
		g.Route(`GET`, ``, Index).SetName(`user.wallet`)
		g.Route(`GET,POST`, `/flow`, Flow).SetName(`user.wallet.flow`)
		g.Route(`GET,POST`, `/recharge`, Recharge).SetName(`user.wallet.recharge`)
		g.Route(`POST`, `/recharge/prepaid_card`, PrepaidCard).SetName(`user.wallet.prepaid_card`) //使用充值卡充值
	})
}
