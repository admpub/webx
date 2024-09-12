package image

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/image`, func(g echo.RouteRegister) {
		g.Route(`GET`, `/proxy/:token`, Proxy)
	})
}
