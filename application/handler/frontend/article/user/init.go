package user

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		g := u.Group(`/article`)
		g.Route(`GET,POST`, `/create`, Create).SetName(`user.article.create`)
		g.Route(`GET,POST`, `/edit/:id`, Edit).SetName(`user.article.edit`)
		g.Get(`/list`, List).SetName(`user.article.list`)
		g.Route(`GET,POST`, `/delete/:id`, Delete).SetName(`user.article.delete`)
	})
}
