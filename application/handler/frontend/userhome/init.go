package userhome

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/u`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/:customerId`, Index).SetName(`user.home`)
		g.Route(`GET,POST`, `/:customerId/:operate`, Index).SetName(`user.home.page`)
	})
}
