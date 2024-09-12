package user

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		g := u.Group(`/article`)
		g.Route(`GET,POST`, `/create`, Create)
		g.Route(`GET,POST`, `/edit/:id`, Edit)
		g.Get(`/list`, List)
		g.Route(`GET,POST`, `/delete/:id`, Delete)
	})
}
