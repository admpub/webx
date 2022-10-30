package user

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/registry/dashboard"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/user`, func(g echo.RouteRegister) {
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
