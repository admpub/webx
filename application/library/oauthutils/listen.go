package oauthutils

import (
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo/handler/oauth2"
	"github.com/webx-top/echo/subdomains"
)

func init() {
	config.OnKeySetSettings(`base.siteURL`, onChangeFrontendURL)
}

func onChangeFrontendURL(d config.Diff) error {
	if defaultOAuth == nil || !d.IsDiff {
		return nil
	}
	host := d.String()
	if len(host) == 0 {
		host = subdomains.Default.URL(``, `frontend`)
	}
	defaultOAuth.HostURL = host
	defaultOAuth.Config.RangeAccounts(func(account *oauth2.Account) bool {
		// 清空生成的网址，以便于在后面的 GenerateProviders() 函数中重新生成新的网址
		account.CallbackURL = ``
		account.LoginURL = ``
		return true
	})
	defaultOAuth.Config.GenerateProviders()
	return nil
}
