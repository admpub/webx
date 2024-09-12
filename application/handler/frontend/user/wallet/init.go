package wallet

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 钱包
		g := u.Group(`/wallet`)
		g.Route(`GET`, ``, Index)
		g.Route(`GET,POST`, `/flow`, Flow)
		g.Route(`GET,POST`, `/recharge`, Recharge)
		g.Route(`POST`, `/recharge/prepaid_card`, PrepaidCard) //使用充值卡充值
	})
}
