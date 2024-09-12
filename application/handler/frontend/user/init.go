package user

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/agent"
	_ "github.com/admpub/webx/application/handler/frontend/user/membership"
	_ "github.com/admpub/webx/application/handler/frontend/user/profile"
	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webfront/initialize/frontend"
	xMW "github.com/coscms/webfront/middleware"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 用户个人
		u.Route(`GET`, `/index`, Index)
		u.Route(`GET,POST`, `/message/unread_count`, MessageUnreadCount)
		u.Route(`GET,POST`, `/message/inbox`, MessageInbox)
		u.Route(`GET,POST`, `/message/outbox`, MessageOutbox)
		u.Route(`GET,POST`, `/message/system`, SystemMessage)
		u.Route(`GET,POST`, `/message/view/:type/:id`, MessageView)
		u.Route(`GET,POST`, `/message/send`, MessageSendHandler)

		g := u.Group(`/file`)
		// 上传
		g.Route(`POST`, `/upload`, Upload)
		// 裁剪图片
		g.Route(`GET,POST`, `/crop`, Crop)
		// 图片管理
		g.Route(`GET,POST`, `/finder`, Finder)

	}, xMW.AuthCheck)

}
