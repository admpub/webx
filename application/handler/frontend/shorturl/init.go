package shorturl

import (
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/r`, func(u echo.RouteRegister) {
		u.Route(`GET,POST`, `/:shortId`, Find).SetName(`short_url`)
	})
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		shortURLG := u.Group(`/short_url`)
		shortURLG.Route(`GET,POST`, `/create`, Create).SetName(`user.short_url.create`)
		shortURLG.Route(`GET,POST`, `/edit/:id`, Edit).SetName(`user.short_url.edit`)
		shortURLG.Get(`/list`, List).SetName(`user.short_url.list`)
		shortURLG.Route(`GET,POST`, `/delete/:id`, Delete).SetName(`user.short_url.delete`)
		shortURLG.Route(`GET,POST`, `/analysis`, Analysis).SetName(`user.short_url.analysis`)
	})
}
