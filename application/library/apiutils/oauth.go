package apiutils

import (
	"fmt"
	"net/url"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"

	"github.com/admpub/webx/application/library/oauthutils"
)

type OauthProvider struct {
	Name      string `json:"name" xml:"name"`
	English   string `json:"english" xml:"english"`
	IconClass string `json:"iconClass" xml:"iconClass"`
	IconImage string `json:"iconImage" xml:"iconImage"`
	WrapClass string `json:"wrapClass" xml:"wrapClass"`
	LoginURL  string `json:"loginURL" xml:"loginURL"`
}

func OauthProvidersFrom(accounts []oauth2.Account) []*OauthProvider {
	var providers []*OauthProvider
	for _, item := range accounts {
		if !item.On {
			continue
		}
		provider := &OauthProvider{
			Name:      item.Name,
			English:   item.Name,
			IconClass: ``,
			IconImage: ``,
			LoginURL:  item.LoginURL,
		}
		if item.Extra != nil {
			provider.IconImage = item.Extra.String(`iconImage`)
			provider.IconClass = item.Extra.String(`iconClass`)
			provider.WrapClass = item.Extra.String(`wrapClass`)
			title := item.Extra.String(`title`)
			if len(title) > 0 {
				provider.Name = title
			}
		}
		providers = append(providers, provider)
	}
	return providers
}

type OauthProvidersResponse struct {
	List []*OauthProvider `json:"list"`
}

type OauthOption interface {
	GetAccountID() uint64
	ApplySetting() (err error)
	GetAppID() string
	OauthProvierListURL(appID string) (string, error)
}

var OauthOptionsCreater = func(ctx echo.Context, typ Type, generators ...URLValuesGenerator) OauthOption {
	return NewOptions(ctx, TypeOauth, generators...)
}

func OauthProviders(ctx echo.Context) ([]*OauthProvider, error) {
	apiOpt := OauthOptionsCreater(ctx, TypeOauth)
	accountID := apiOpt.GetAccountID()
	if accountID <= 0 {
		return OauthProvidersFrom(oauthutils.Accounts()), nil
	}
	if err := apiOpt.ApplySetting(); err != nil {
		return nil, err
	}
	appID := apiOpt.GetAppID()
	if len(appID) == 0 {
		return OauthProvidersFrom(oauthutils.Accounts()), nil
	}
	apiURL, err := apiOpt.OauthProvierListURL(appID)
	if err != nil {
		return nil, err
	}
	platformList := &OauthProvidersResponse{}
	data := echo.NewData(ctx)
	data.Data = platformList
	_, err = SubmitWithRecv(ctx, data, apiURL, url.Values{})
	if err == nil && data.Code.Int() != 1 {
		err = fmt.Errorf(`OauthProviders: %v`, data.Info)
	}
	return platformList.List, err
}

func GetOauthProviderTitle(list []*OauthProvider, name string) string {
	for _, v := range list {
		if v.English == name {
			return v.Name
		}
	}
	return ``
}

func GetOauthProvider(list []*OauthProvider, name string) *OauthProvider {
	for _, v := range list {
		if v.English == name {
			return v
		}
	}
	return nil
}

// OauthProvierListURL 社区登录供应商列表
func (o *Options) OauthProvierListURL(appID string) (string, error) {
	urlValues := url.Values{}
	urlValues.Set(`appID`, appID)
	o.SetGenerator(DefaultURLValuesGenerator(urlValues))
	uri, formData, err := o.ToURL(`/open/v1/oauth/providers`)
	if err != nil {
		return ``, err
	}
	return uri + `?` + formData.Encode(), nil
}
