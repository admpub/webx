package user

import (
	"github.com/coscms/webcore/registry/dashboard"
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/user`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/message/unread_count`, MessageUnreadCount)
		g.Route(`GET,POST`, `/message/inbox`, MessageInbox)
		g.Route(`GET,POST`, `/message/outbox`, MessageOutbox)
		g.Route(`GET,POST`, `/message/system`, SystemMessage)
		g.Route(`GET,POST`, `/message/delete`, MessageDelete)
		g.Route(`GET,POST`, `/message/view/:type/:id`, MessageView)
		g.Route(`GET,POST`, `/message/send`, MessageSendHandler)
	})

	dashboard.TopButtonRegister(&dashboard.Button{
		Tmpl: `official/user/topbutton/message`,
	})
	dashboard.TopButtonRegister(&dashboard.Button{
		Tmpl: `official/user/topbutton/notice`,
	})
	dashboard.TopButtonRegister(&dashboard.Button{
		Tmpl: `official/user/topbutton/home`,
	})
	dashboard.GlobalFooterRegister(&dashboard.GlobalFooter{
		Tmpl: `official/user/footer/footer`,
	})
}
