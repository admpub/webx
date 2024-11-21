package user

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/membership"
	_ "github.com/admpub/webx/application/handler/frontend/user/profile"
	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/xnotice"
	xMW "github.com/coscms/webfront/middleware"
)

func init() {
	xnotice.RegsterCmder()
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		xnotice.RegisterRoute(u)
		// 用户个人
		u.Route(`GET`, `/index`, Index).SetName(`user.index`)
		u.Route(`GET,POST`, `/message/unread_count`, MessageUnreadCount).SetName(`user.message.unread_count`)
		u.Route(`GET,POST`, `/message/inbox`, MessageInbox).SetName(`user.message.inbox`)
		u.Route(`GET,POST`, `/message/outbox`, MessageOutbox).SetName(`user.message.outbox`)
		u.Route(`GET,POST`, `/message/system`, SystemMessage).SetName(`user.message.system`)
		u.Route(`GET,POST`, `/message/view/:type/:id`, MessageView).SetName(`user.message.view`)
		u.Route(`GET,POST`, `/message/send`, MessageSendHandler).SetName(`user.message.send`)

		// 个人收藏夹
		favoriteG := u.Group(`/favorite`)
		favoriteG.Route(`GET`, `/index`, favoriteList).SetName(`user.favorite`)
		favoriteG.Route(`GET,POST`, `/delete`, favoriteDelete).SetName(`user.favorite.delete`)
		favoriteG.Route(`GET,POST`, `/go/:id`, favoriteGo).SetName(`user.favorite.go`)

		// 个人文件
		g := u.Group(`/file`)
		// 上传
		g.Route(`POST`, `/upload`, Upload).SetName(`user.file.upload`)
		// 裁剪图片
		g.Route(`GET,POST`, `/crop`, Crop).SetName(`user.file.crop`)
		// 图片管理
		g.Route(`GET,POST`, `/finder`, Finder).SetName(`user.file.finder`)

	}, xMW.AuthCheck)

	bootconfig.OnStart(-1, func() {
		if !config.IsInstalled() {
			return
		}
		if !config.FromFile().Extend.Bool(`disableAutoResetClientCount`) {
			go xnotice.ResetClientCount()
		}
	})
}
