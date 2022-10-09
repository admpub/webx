package shorturl

import (
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.RegisterToGroup(`/r`, func(u echo.RouteRegister) {
		u.Route(`GET,POST`, `/:shortId`, Find)
	})
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		shortURLG := u.Group(`/short_url`)
		shortURLG.Route(`GET,POST`, `/create`, Create)
		shortURLG.Route(`GET,POST`, `/edit/:id`, Edit)
		shortURLG.Get(`/list`, List)
		shortURLG.Route(`GET,POST`, `/delete/:id`, Delete)
		shortURLG.Route(`GET,POST`, `/analysis`, Analysis)
	})
}
