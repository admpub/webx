package providers

import (
	"github.com/markbates/goth"

	"github.com/webx-top/echo/handler/oauth2"

	"github.com/coscms/oauth2s/client/goth/providers/alipay"
	"github.com/coscms/oauth2s/client/goth/providers/qq"
	"github.com/coscms/oauth2s/client/goth/providers/wechat"
	"github.com/coscms/oauth2s/client/goth/providers/weibo"
)

var constructors = map[string]func(account *oauth2.Account) goth.Provider{
	`alipay`: func(account *oauth2.Account) goth.Provider {
		return alipay.New(account.Key, account.Secret, account.CallbackURL, true)
	},
	`alipay_dev`: func(account *oauth2.Account) goth.Provider {
		return alipay.New(account.Key, account.Secret, account.CallbackURL, false)
	},
	`qq`: func(account *oauth2.Account) goth.Provider {
		return qq.New(account.Key, account.Secret, account.CallbackURL)
	},
	`weibo`: func(account *oauth2.Account) goth.Provider {
		return weibo.New(account.Key, account.Secret, account.CallbackURL)
	},
	`wechat`: func(account *oauth2.Account) goth.Provider {
		return wechat.New(account.Key, account.Secret, account.CallbackURL)
	},
}

func Register(name string, constructor func(account *oauth2.Account) goth.Provider) {
	constructors[name] = constructor
}

func Get(name string) func(account *oauth2.Account) goth.Provider {
	constructor, _ := constructors[name]
	return constructor
}
