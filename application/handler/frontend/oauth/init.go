package oauth

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/apiutils"
	xMW "github.com/coscms/webfront/middleware"
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
		g3 := oauthG.Group(`/other`)
		wechatG := g3.Group(`/wechat`) // 微信 /oauth/other/wechat
		{
			wechatG.Route(`GET,POST`, `/mp`, WechatMP) // 微信小程序 /oauth/other/wechat/mp

			//TODO:扫码关注公众号登录
			//wechatG.Route(`GET,POST`, `/gh`, WechatGH)                  // 微信公众号 /oauth/other/wechat/gh
			//wechatG.Get(`/gh/callback`, WechatGHCheckSign) // 微信公众号 /oauth/other/wechat/gh/callback
			//wechatG.Post(`/gh/callback`, WechatGHCallback) // 微信公众号 /oauth/other/wechat/gh/callback
		}
	})
	backend.OnInstalled(onInstalled)
}
