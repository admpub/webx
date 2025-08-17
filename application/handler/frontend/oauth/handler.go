package oauth

import (
	"github.com/admpub/goth"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/handler/oauth2"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webfront/library/oauthutils"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

// 通过oauth登录第三方网站成功之后的处理
func successHandler(ctx echo.Context) error {
	if config.Setting(`base`).String(`customerLogin`, `open`) == `close` {
		return ctx.E(`本站已经暂时关闭登录，请稍后再尝试`)
	}
	oauthM := modelCustomer.NewOAuth(ctx)
	var (
		ouser       *goth.User
		fromSession bool
	)
	if user := Default().User(ctx); len(user.Provider) > 0 {
		ouser = &user
	}
	if ouser == nil { //几乎不会执行到
		var err error
		ouser, _, err = oauthM.GetSession()
		if err != nil {
			return err
		}
		if ouser == nil {
			return ctx.E(`第三方账号登录失败`)
		}
		fromSession = true
		defer func() {
			if fromSession {
				oauthM.DelSession()
			}
		}()
	}
	if len(ouser.UserID) == 0 {
		return ctx.NewError(code.InvalidParameter, `oauth2登录后获取UserID无效`)
	}
	end, err := oauthutils.FireAfterLoginSuccess(ctx, ouser)
	if err != nil || end {
		return err
	}
	next := ctx.Form(echo.DefaultNextURLVarName)
	var needSignIn bool
	_, err = checkOrUpdateUser(ctx, oauthM, ouser, func(ctx echo.Context) (bool, error) {
		if !ctx.Queryx(`force-create`).Bool() { //绑定用户账号，除非强制指定了自动创建新账号
			if !fromSession { // 如果gothUser数据不是来自session，则需要保存到session中以便于用户登录或注册账号后关联oauth数据
				oauthM.SaveSession(ouser)
			}
			fromSession = false
			needSignIn = true
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	if len(next) == 0 {
		next = common.GetSavedNextURL(ctx, ctx.URLFor(`/index`))
	}
	if needSignIn {
		next = ctx.URLFor(`/sign_in?next=` + com.URLEncode(next))
	}
	return ctx.Redirect(next)
}

// checkOrUpdateUser 检查或更新用户信息
func checkOrUpdateUser(ctx echo.Context, oauthM *modelCustomer.OAuth, ouser *goth.User, newUserSignupBefore func(ctx echo.Context) (bool, error)) (*modelCustomer.Customer, error) {
	var customerM *modelCustomer.Customer
	customer := sessdata.Customer(ctx)
	err := oauthM.GetByOutUser(ouser)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return customerM, err
		}

		// === 未曾绑定过 ====================

		// ==============================
		// 操作客户资料
		// ==============================
		if customer == nil { //未登录
			if newUserSignupBefore != nil {
				exit, err := newUserSignupBefore(ctx)
				if err != nil {
					return customerM, err
				}
				if exit {
					return customerM, nil
				}
			}
			// 用户不存在，需要新建并自动登录
			// 注册成功会自动登录
			customerM, err = oauthM.SignUpCustomer(ouser)
			if err != nil {
				return customerM, err
			}
			if customerM.Id < 1 {
				return customerM, ctx.E(`注册时发生异常：获取InsertID失败`)
			}
			customer = customerM.OfficialCustomer
		} else { //已登录时，看看是否需要更新未填写的资料
			customer.SetContext(ctx)
			customerM = modelCustomer.NewCustomer(ctx)
			customerM.OfficialCustomer = customer
			set := echo.H{}
			if len(customer.Email) < 1 && len(ouser.Email) > 0 {
				set[`email`] = ouser.Email
			}
			if len(customer.Avatar) < 1 && len(ouser.AvatarURL) > 0 {
				set[`avatar`] = ouser.AvatarURL
			}
			gender := oauthM.OAuthUserGender(ouser)
			if len(gender) > 0 && (len(customer.Gender) < 1 || (customer.Gender == `secret` && gender != customer.Gender)) {
				set[`gender`] = gender
			}
			if len(set) > 0 {
				err := customerM.UpdateFields(nil, set, `id`, customer.Id)
				if err != nil {
					log.Error(ctx.T(`更新 official_customer 表数据行#%d为 %s 时，发生错误: `, customer.Id, com.Dump(set, false)), err)
				}
			}
		}

		// ==============================
		// 添加OAuth关联记录
		// ==============================
		oauthM.CustomerId = customer.Id
		oauthM.CopyFrom(ouser)
		_, err = oauthM.Add()
		return customerM, err
	}

	// === 已经绑定过 ====================

	if customer != nil && oauthM.CustomerId != customer.Id {
		err = ctx.NewError(code.DataUnavailable, `此外部账号已经被其他用户绑定`)
		return customerM, err
	}

	oauthSet := echo.H{}
	if ouser.AccessToken != oauthM.AccessToken {
		oauthSet[`access_token`] = ouser.AccessToken
	}
	if ouser.RefreshToken != oauthM.RefreshToken {
		oauthSet[`refresh_token`] = ouser.RefreshToken
	}
	if !ouser.ExpiresAt.IsZero() && ouser.ExpiresAt.Unix() > 0 {
		oauthSet[`expired`] = ouser.ExpiresAt.Unix()
	}

	// 直接登录
	customerM = modelCustomer.NewCustomer(ctx)
	err = customerM.Get(nil, `id`, oauthM.CustomerId)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return customerM, err
		}
		// 用户不存在，需要重新注册（有可能用户已经被删除）
		// 注册成功会自动登录
		customerM, err = oauthM.SignUpCustomer(ouser)
		if err != nil {
			return customerM, err
		}
		if customerM.Id < 1 {
			return customerM, ctx.E(`注册时发生异常：获取InsertID失败`)
		}
		oauthSet[`customer_id`] = customerM.Id
		err = oauthM.UpdateFields(nil, oauthSet, db.Cond{`id`: oauthM.Id})
		return customerM, err
	}
	// 更新用户的旧资料
	if len(ouser.AvatarURL) > 0 && ouser.AvatarURL != oauthM.Avatar {
		oauthSet[`avatar`] = ouser.AvatarURL
		if oauthM.Avatar == customerM.Avatar {
			err = customerM.UpdateField(nil, `avatar`, ouser.AvatarURL, `id`, customerM.Id)
			if err != nil {
				log.Error(ctx.T(`更新本地用户头像为oauth2用户头像时失败`), `: `, err.Error())
			}
		}
	}

	if len(oauthSet) > 0 {
		err = oauthM.UpdateFields(nil, oauthSet, `id`, oauthM.Id)
		if err != nil {
			log.Error(ctx.T(`更新用户oauth2的数据(%s)失败`, echo.Dump(oauthSet, false)), `: `, err.Error())
		}
	}

	// 未登录时设置登录状态
	if customer == nil {
		co := modelCustomer.NewCustomerOptions(customerM.OfficialCustomer)
		co.SignInType = `oauth2.` + ouser.Provider
		co.ApplyOptions(modelCustomer.GenerateOptionsFromHeader(ctx)...)
		err = customerM.FireSignInSuccess(co, model.AuthTypeOauth2)
	}

	return customerM, err
}

func beginAuthHandler(ctx echo.Context) error {
	return oauth2.BeginAuthHandler(ctx)
}
