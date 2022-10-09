package apiutils

import (
	"strings"

	"github.com/admpub/nging/v4/application/library/config"
	webxdbschema "github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/xcommon"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func (o *Options) getAccount(cond db.Compound) (err error) {
	accountM := webxdbschema.NewOfficialCommonApiAccount(o.ctx)
	err = accountM.Get(nil, cond)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = o.ctx.NewError(code.DataNotFound, `没找到API账号数据：%s`, echo.Dump(cond, false))
		}
		return
	}
	o.Account = accountM
	if config.FromFile().Sys.IsEnv(`prod`) {
		o.URLPrefix = accountM.Url
	} else {
		o.URLPrefix = accountM.UrlDev
	}
	o.URLPrefix = strings.TrimSuffix(o.URLPrefix, `/`)
	return
}

var AppInfoDefaultGetter = func(ctx echo.Context, cond db.Compound) (appInfo AppInfo, err error) {
	err = ctx.NewError(code.Unsupported, `尚未设置App信息获取方式`)
	return
}

func (o *Options) getApp(cond db.Cond) (err error) {
	var appInfo AppInfo
	if o.appInfoGetter == nil {
		appInfo, err = AppInfoDefaultGetter(o.ctx, cond)
	} else {
		appInfo, err = o.appInfoGetter(o.ctx, cond)
	}
	if err != nil {
		return
	}
	o.App = appInfo
	o.URLPrefix = xcommon.SiteURL(o.ctx)
	return
}
