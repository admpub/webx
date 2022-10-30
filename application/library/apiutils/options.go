package apiutils

import (
	"net/url"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/null"
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewOptions(ctx echo.Context, typ Type, generators ...URLValuesGenerator) *Options {
	var generator URLValuesGenerator
	if len(generators) > 0 {
		generator = generators[0]
	}
	return &Options{
		ctx:       ctx,
		generator: generator,
		Type:      typ,
	}
}

type AppInfo interface {
	GetAppSecret() string
	GetAppId() string
}

type Options struct {
	ctx           echo.Context
	generator     URLValuesGenerator
	signaturer    func(url.Values) string
	appInfoGetter func(ctx echo.Context, cond db.Cond) (appInfo AppInfo, err error)
	applied       bool
	accountID     null.Uint64
	Account       *dbschema.OfficialCommonApiAccount
	App           AppInfo
	URLPrefix     string
	Type          Type
}

func (o *Options) SetGenerator(g URLValuesGenerator) *Options {
	o.generator = g
	return o
}

func (o *Options) SetSignaturer(fn func(url.Values) string) *Options {
	o.signaturer = fn
	return o
}

func (o *Options) SetAppInfoGetter(appInfoGetter func(ctx echo.Context, cond db.Cond) (appInfo AppInfo, err error)) *Options {
	o.appInfoGetter = appInfoGetter
	return o
}

func (o *Options) GetAccountID() uint64 {
	if o.accountID.Valid {
		return o.accountID.Uint64
	}
	o.accountID.Valid = true
	o.accountID.Uint64 = config.Setting(`thirdparty`).Uint64(string(o.Type))
	return o.accountID.Uint64
}

func (o *Options) ApplySetting() (err error) {
	if o.applied {
		return
	}
	o.applied = true
	accountID := o.GetAccountID()
	if accountID > 0 {
		err = o.getAccount(db.Cond{`id`: accountID})
	} else {
		err = o.getApp(db.Cond{`owner_type`: `official`})
	}
	return
}

func (o *Options) GetAppID() string {
	if o.Account != nil {
		return o.Account.AppId
	}
	if o.App != nil {
		return o.App.GetAppId()
	}
	return ``
}

func (o *Options) GetAppSecret() string {
	if o.Account != nil {
		return o.Account.AppSecret
	}
	if o.App != nil {
		return o.App.GetAppSecret()
	}
	return ``
}

func (o *Options) Context() echo.Context {
	return o.ctx
}
