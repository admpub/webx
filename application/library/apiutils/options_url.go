package apiutils

import (
	"net/url"

	"github.com/admpub/webx/application/library/top"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/code"
)

func (o *Options) ToURL(urlPath string, strength ...bool) (uri string, formData url.Values, err error) {
	formData = o.generator.URLValues(o.signaturer)
	appID := formData.Get(`appID`)
	if len(appID) > 0 {
		oldAppID := o.GetAppID()
		if len(oldAppID) > 0 {
			if oldAppID != appID {
				err = o.ctx.NewError(code.DataUnavailable, `AppID前后不一致`)
			}
		} else {
			if o.GetAccountID() > 0 {
				err = o.getAccount(db.Cond{`app_id`: appID})
			} else {
				err = o.getApp(db.Cond{`app_id`: appID})
			}
		}
	} else {
		err = o.ApplySetting()
	}
	if err != nil {
		return
	}
	var appSecret string
	if o.Account != nil {
		appSecret = o.Account.AppSecret
		if len(appID) == 0 {
			formData.Set(`appID`, o.Account.AppId)
		}
	} else {
		appSecret = o.App.GetAppSecret()
		if len(appID) == 0 {
			formData.Set(`appID`, o.App.GetAppId())
		}
	}
	if len(strength) > 0 && strength[0] {
		appSecret = top.StrengthenSafeSecret(o.ctx, appSecret) // 加强防篡改安全性
	}
	formData = BuildURLValues(formData, appSecret)
	uri = o.URLPrefix + urlPath
	return
}
