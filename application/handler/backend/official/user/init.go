package user

import (
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/webx-top/echo"
)

func init() {
	httpserver.Backend.Router.RegisterToGroup(`/user`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/message/unread_count`, MessageUnreadCount)
		g.Route(`GET,POST`, `/message/inbox`, MessageInbox)
		g.Route(`GET,POST`, `/message/outbox`, MessageOutbox)
		g.Route(`GET,POST`, `/message/system`, SystemMessage)
		g.Route(`GET,POST`, `/message/delete`, MessageDelete)
		g.Route(`GET,POST`, `/message/view/:type/:id`, MessageView)
		g.Route(`GET,POST`, `/message/send`, MessageSendHandler)
	})

	httpserver.Backend.Dashboard.TopButtons.Register(&dashboard.Button{
		Tmpl: `official/user/topbutton/message`,
	})
	httpserver.Backend.Dashboard.TopButtons.Register(&dashboard.Button{
		Tmpl: `official/user/topbutton/notice`,
	})
	httpserver.Backend.Dashboard.TopButtons.Register(&dashboard.Button{
		Tmpl: `official/user/topbutton/home`,
	})
	httpserver.Backend.Dashboard.GlobalFooters.Register(&dashboard.GlobalFooter{
		Tmpl: `official/user/footer/footer`,
	})
}
