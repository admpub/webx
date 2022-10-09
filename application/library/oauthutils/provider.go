package oauthutils

import (
	"github.com/coscms/oauth2s/client/goth/providers"
	"github.com/markbates/goth"
	"github.com/webx-top/echo/handler/oauth2"

	// - oauth2 provider
	//"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/microsoftonline"
)

// RegisterProvider 注册Provider
func RegisterProvider(c *oauth2.Config) {

	providers.Register(`microsoft`, func(account *oauth2.Account) goth.Provider {
		if len(account.CallbackURL) == 0 {
			account.CallbackURL = c.CallbackURL(account.Name)
		}
		m := microsoftonline.New(account.Key, account.Secret, account.CallbackURL)
		m.SetName(`microsoft`)
		return m
	})

	/*
		providers.Register(`apple`, func(account *oauth2.Account) goth.Provider {
			if len(account.CallbackURL) == 0 {
				account.CallbackURL = c.CallbackURL(account.Name)
			}
			return apple.New(account.Key, account.Secret, account.CallbackURL)
		})
	*/
}
