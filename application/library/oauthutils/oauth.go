package oauthutils

import (
	"github.com/admpub/log"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/markbates/goth"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/subdomains"

	dbschemaNging "github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/config"
	oauthLibrary "github.com/admpub/nging/v5/application/library/oauth"
	"github.com/admpub/nging/v5/application/registry/settings"
	"github.com/coscms/oauth2s/client/goth/providers"
)

var (
	defaultOAuth      *oauth2.OAuth
	SuccessHandler    interface{}
	BeginAuthHandler  echo.Handler
	AfterLoginSuccess []func(ctx echo.Context, ouser *goth.User) (end bool, err error)
)

func OnAfterLoginSuccess(hooks ...func(ctx echo.Context, ouser *goth.User) (end bool, err error)) {
	AfterLoginSuccess = append(AfterLoginSuccess, hooks...)
}

func FireAfterLoginSuccess(ctx echo.Context, ouser *goth.User) (end bool, err error) {
	for _, hook := range AfterLoginSuccess {
		end, err = hook(ctx, ouser)
		if err != nil || end {
			return
		}
	}
	return
}

func Default() *oauth2.OAuth {
	return defaultOAuth
}

// InitOauth 第三方登录
func InitOauth(e *echo.Echo) {
	if config.IsInstalled() {
		settings.Init(nil)
	}
	host := subdomains.Default.URL(``, frontend.Name)
	oauth2Config := &oauth2.Config{}
	RegisterProvider(oauth2Config)

	if config.IsInstalled() {
		if oauthAccounts, err := FindAccounts(); err != nil {
			log.Error(err)
		} else {
			oauth2Config.AddAccount(oauthAccounts...)
		}
	}

	defaultOAuth = oauth2.New(host, oauth2Config)
	defaultOAuth.SetSuccessHandler(SuccessHandler)
	defaultOAuth.SetBeginAuthHandler(BeginAuthHandler)
	defaultOAuth.Wrapper(e, session.Middleware(config.SessionOptions))
}

func Accounts() []*oauth2.Account {
	return Default().Config.Accounts
}

// FindAccounts 第三方登录平台账号
func FindAccounts() ([]*oauth2.Account, error) {
	m := &dbschemaNging.NgingConfig{}
	_, err := m.ListByOffset(nil, nil, 0, -1, db.Cond{`group`: `oauth`})
	var result []*oauth2.Account
	isProduction := config.FromFile().Sys.IsEnv(`prod`)
	for _, row := range m.Objects() {
		if len(row.Value) == 0 {
			continue
		}
		cfg := &oauthLibrary.Config{Name: row.Key, On: row.Disabled != `Y`}
		err = com.JSONDecode([]byte(row.Value), cfg)
		if err != nil {
			return result, err
		}
		value := cfg.ToAccount(row.Key)
		var provider func(account *oauth2.Account) goth.Provider
		if !isProduction {
			provider = providers.Get(value.Name + `_dev`)
		}
		if provider == nil {
			provider = providers.Get(value.Name)
		}
		value.SetConstructor(provider)
		if value.On {
			result = append(result, value)
		}
	}
	return result, err
}

// UpdateAccount 第三方登录平台账号
func UpdateAccount() error {
	accounts, err := FindAccounts()
	if err != nil {
		return err
	}
	Default().Config.Accounts = accounts
	Default().Config.GenerateProviders()
	log.Debug(`Update oauth configuration information`)
	return nil
}
