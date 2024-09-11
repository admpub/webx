package oauth

import (
	"github.com/coscms/webcore/library/config"
	indexHanlder "github.com/admpub/webx/application/handler/frontend/index"
	wechatmp "github.com/admpub/webx/application/library/thirdparty/login/wechat/mp"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
)

// WechatMP 微信小程序登录
func WechatMP(ctx echo.Context) error {
	// 已经登录的时候跳过当前页面
	if sessdata.Customer(ctx) != nil {
		ctx.Data().SetInfo(ctx.T(`已经登录过了`))
		next := ctx.Form(`next`)
		next = echo.GetOtherURL(ctx, next)
		if len(next) == 0 {
			next = sessdata.URLFor(`/index`)
		}
		return ctx.Redirect(next)
	}
	data := ctx.Data()
	cfg := config.Setting(`oauth`, `wechat`)
	if len(cfg) == 0 {
		err := ctx.NewError(stdCode.Unsupported, ctx.T(`不支持此种登录方式`))
		return ctx.JSON(data.SetError(err))
	}
	post := wechatmp.NewWechatPostData()
	if err := ctx.MustBind(post); err != nil {
		return ctx.JSON(data.SetError(err))
	}
	if err := post.Check(ctx); err != nil {
		return ctx.JSON(data.SetError(err))
	}
	appID := cfg.String(`key`)
	appSecret := cfg.String(`secret`)
	result, err := post.Post(ctx, appID, appSecret)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	ouser := result.AsUser(post)
	oauthM := modelCustomer.NewOAuth(ctx)
	var customerM *modelCustomer.Customer
	customerM, err = checkOrUpdateUser(ctx, oauthM, ouser, nil)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}

	err = indexHanlder.SetJWTData(ctx, customerM)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	return ctx.JSON(data.SetInfo(ctx.T(`登录成功`)))
}
