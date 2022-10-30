package oauth

import (
	"github.com/admpub/nging/v5/application/handler/setup"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/apiutils"
	xMW "github.com/admpub/webx/application/middleware"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	xMW.Use(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`OAuthAccounts`, func() []*apiutils.OauthProvider {
				list, _ := apiutils.OauthProviders(c)
				return list
			})
			return h.Handle(c)
		})
	})
	frontend.Register(func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/oauth_sign_in`, SignIn)
		initOauth(frontend.IRegister().Echo())

		oauthG := g.Group(oauth2.DefaultPath)
		mpG := oauthG.Group(`/mp`)                 // 小程序
		mpG.Route(`GET,POST`, `/wechat`, MPWechat) // 微信小程序
	})
	setup.OnInstalled(onInstalled)
}
