package profile

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/handler/frontend/user/binding"
	"github.com/admpub/webx/application/library/apiutils"
	"github.com/admpub/webx/application/library/xcommon"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

// 跳转到第三方平台进行登录验证
func bindingOAuth(ctx echo.Context, customer *modelCustomer.Customer) error {
	provider := ctx.Formx(`provider`).String()
	cancelID := ctx.Formx(`cancelId`).Uint64()
	if cancelID > 0 {
		return unbindOAuth(ctx, customer, provider, cancelID)
	}
	siteURL := xcommon.SiteURL(ctx)
	nextURL := com.URLEncode(siteURL + `/user/profile/binding?type=oauth`)
	providers, err := apiutils.OauthProviders(ctx)
	if err != nil {
		return err
	}
	item := apiutils.GetOauthProvider(providers, provider)
	if item == nil {
		return ctx.NewError(code.Unsupported, `不支持绑定: %v`, provider)
	}
	loginURL := com.WithURLParams(item.LoginURL, `next`, nextURL)
	return ctx.Redirect(loginURL)
}

func unbindOAuth(ctx echo.Context, customer *modelCustomer.Customer, provider string, cancelID uint64) error {
	oauthM := modelCustomer.NewOAuth(ctx)
	var cancelable bool
	var err error
	if customer.RegisteredBy == `oauth2.login` {
		var exists bool
		exists, err = oauthM.ExistsOtherBinding(customer.Id, cancelID)
		if err != nil {
			return err
		}
		cancelable = exists
		if !exists {
			err = ctx.NewError(code.Failure, `当前账号是用快捷登录进行注册，不能全部取消绑定，否则将导致无法登录`)
		}
	} else {
		cancelable = true
	}
	if !cancelable {
		return err
	}
	err = oauthM.Delete(nil, db.And(
		db.Cond{`id`: cancelID},
		db.Cond{`customer_id`: customer.Id},
	))
	if err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(xMW.URLFor(`/user/profile/binding?type=oauth`))
}

// 查询支持绑定的第三方平台账号和当前账号是否已与其绑定
func bindingOAuthGet(ctx echo.Context, m *modelCustomer.Customer) error {
	oauthM := modelCustomer.NewOAuth(ctx)
	_, err := oauthM.ListByOffset(nil, nil, 0, 200, db.And(
		db.Cond{`customer_id`: m.Id},
	))
	if err != nil {
		return err
	}
	bindedOAuthAccounts := map[string][]*dbschema.OfficialCustomerOauth{}
	for _, row := range oauthM.Objects() {
		if _, ok := bindedOAuthAccounts[row.Type]; !ok {
			bindedOAuthAccounts[row.Type] = []*dbschema.OfficialCustomerOauth{}
		}
		bindedOAuthAccounts[row.Type] = append(bindedOAuthAccounts[row.Type], row)
	}
	providers, err := apiutils.OauthProviders(ctx)
	if err != nil {
		return err
	}
	oAuthProviders := make([]*binding.BindOAuthAccount, len(providers))
	for index, item := range providers {
		provider := item.English
		oAuthTitle := item.Name
		users, binded := bindedOAuthAccounts[provider]
		oAuthProviders[index] = &binding.BindOAuthAccount{
			IconClass:   item.IconClass,
			IconImage:   item.IconImage,
			WrapClass:   item.WrapClass,
			Provider:    provider,
			Title:       oAuthTitle,
			Description: ctx.T(`将%s与当前账号绑定`, oAuthTitle),
			Binded:      binded,
			OAuthUsers:  users,
		}
	}
	ctx.Set(`oAuthProviders`, oAuthProviders)
	return nil
}
