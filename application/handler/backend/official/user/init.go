package user

import (
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/webx-top/echo"
)

func init() {
	httpserver.Backend.Router.RegisterToGroup(`/user`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/message/unread_count`, MessageUnreadCount).SetName(`admin.message.unreadCount`)
		g.Route(`GET,POST`, `/message/inbox`, MessageInbox).SetName(`admin.message.inbox`)
		g.Route(`GET,POST`, `/message/outbox`, MessageOutbox).SetName(`admin.message.outbox`)
		g.Route(`GET,POST`, `/message/system`, SystemMessage).SetName(`admin.message.system`)
		g.Route(`GET,POST`, `/message/delete`, MessageDelete).SetName(`admin.message.delete`)
		g.Route(`GET,POST`, `/message/view/:type/:id`, MessageView).SetName(`admin.message.view`)
		g.Route(`GET,POST`, `/message/send`, MessageSendHandler).SetName(`admin.message.send`)
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
