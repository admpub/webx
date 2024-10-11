package oauth

import (
	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/library/oauthutils"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"
)

func init() {
	oauthutils.SuccessHandler = successHandler
	oauthutils.BeginAuthHandler = echo.HandlerFunc(beginAuthHandler)
}

func Default() *oauth2.OAuth {
	return oauthutils.Default()
}

// initOauth 第三方登录
func initOauth(e *echo.Echo) {
	oauthutils.InitOauth(e, httpserver.SearchEngineNoindex())
	if config.IsInstalled() {
		if err := SetSMSConfigs(); err != nil {
			log.Error(err)
		}
	}
}

func Accounts() []oauth2.Account {
	return oauthutils.Accounts()
}

func onInstalled(ctx echo.Context) error {
	if err := SetSMSConfigs(); err != nil {
		log.Error(err)
	}
	return oauthutils.UpdateAccount()
}

// FindAccounts 第三方登录平台账号
func FindAccounts() ([]*oauth2.Account, error) {
	return oauthutils.FindAccounts()
}

// UpdateAccount 第三方登录平台账号
func UpdateAccount() error {
	return oauthutils.UpdateAccount()
}

func SignIn(ctx echo.Context) error {
	return successHandler(ctx)
}
