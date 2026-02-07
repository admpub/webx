package user

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/membership"
	_ "github.com/admpub/webx/application/handler/frontend/user/profile"
	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/xnotice"
	xMW "github.com/coscms/webfront/middleware"
)

func init() {
	xnotice.RegisterCmder()
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		cfg := config.FromFile().Extend.GetStore(`frontendNotice`)
		if !cfg.Bool(`disabled`) {
			xnotice.RegisterRoute(u, cfg)
		}
		// 用户个人
		u.Route(`GET`, `/index`, Index).SetName(`user.index`)

		messageG := u.Group(`/message`)
		messageG.Route(`GET,POST`, `/unread_count`, MessageUnreadCount).SetName(`user.message.unread_count`)
		messageG.Route(`GET,POST`, `/inbox`, MessageInbox).SetName(`user.message.inbox`)
		messageG.Route(`GET,POST`, `/outbox`, MessageOutbox).SetName(`user.message.outbox`)
		messageG.Route(`GET,POST`, `/system`, SystemMessage).SetName(`user.message.system`)
		messageG.Route(`GET,POST`, `/view/:type/:id`, MessageView).SetName(`user.message.view`)
		messageG.Route(`GET,POST`, `/send`, MessageSendHandler).SetName(`user.message.send`)

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

		// 二维码
		qrcodeG := u.Group(`/qrcode`)
		qrcodeG.Route(`GET,POST`, `/scan`, qrcodeScan).SetName(`user.qrcode.scan`).SetMetaKV(httpserver.PermPublicKV())

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
