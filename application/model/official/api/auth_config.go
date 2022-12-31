package api

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/middleware/mwapp"
)

// APIAccountAuthConfig Api账号认证配置
var APIAccountAuthConfig = mwapp.NewAuthConfig().SetSecretGetter(func(ctx echo.Context, appID string) (string, error) {
	m := NewAccount(ctx)
	err := m.Get(nil, `app_id`, appID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `API账号“%s”没找到`, appID).SetZone(`appId`)
		}
		return "", err
	}

	if m.Disabled != `N` {
		return "", ctx.NewError(code.DataUnavailable, `API账号“%s”已经停用`, appID).SetZone(`appId`)
	}
	ctx.Internal().Store(`ApiAccount`, m.OfficialCommonApiAccount)
	return m.AppSecret, nil
})

// APIAuthMiddleware Api账号认证中间件
var APIAuthMiddleware = mwapp.Auth(*APIAccountAuthConfig)
