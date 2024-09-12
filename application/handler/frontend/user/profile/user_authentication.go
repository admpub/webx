package profile

import (
	"github.com/admpub/webx/application/handler/frontend/user/binding"
	"github.com/coscms/webcore/library/common"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
)

func _Authentication(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	m := modelCustomer.NewCustomer(ctx)
	err := m.VerifySession(customer)
	if err != nil {
		if common.IsUserNotLoggedIn(err) {
			return ctx.E(`请先登录`)
		}
		return err
	}
	typ := `password` // 验证类型: 密码验证
	if m.MobileBind == `Y` || m.EmailBind == `Y` {
		typ = ctx.Form(`type`)
		switch typ {
		case `email`, `mobile`:
		default:
			if m.MobileBind == `Y` {
				typ = `mobile`
			} else {
				typ = `email`
			}
		}
	}
	data := ctx.Data()
	ctx.Set(`type`, typ)
	if typ != `password` {
		binder := binding.Get(typ)
		if binder == nil {
			return ctx.E(`不支持的类型: %s`, typ)
		}
		ctx.Set(`objectName`, binder.ObjectName)
		ctx.Set(`typeName`, binder.Name)
		waitingSeconds := BindingSendingInterval(typ)
		maxPerDay := BindingSendingMaxPerDay(typ)
		remainCount := BindingSendingRemainCount(ctx, m, typ)
		ctx.Set(`waitingSeconds`, waitingSeconds)
		ctx.Set(`maxPerDay`, maxPerDay)
		ctx.Set(`remainCount`, remainCount)
	}
	b, err := ctx.Fetch(`user/profile/partial_authentication`, nil)
	if err != nil {
		return err
	}
	ctx.Set(`html`, string(b))
	data.SetData(ctx.Stored())
	return ctx.JSON(data)
}

func modifyPassword(m *modelCustomer.Customer, typ string, newPassword string) error {
	switch typ {
	case `safe`:
		return m.UpdateSafePassword(newPassword)
	default:
		return m.UpdateSignInPassword(newPassword)
	}
}

func checkNewPassword(ctx echo.Context, m *modelCustomer.Customer, typ string) (string, error) {
	switch typ {
	case `safe`:
		safePwd := ctx.Formx(`safePwd`).String()
		safePwd2 := ctx.Formx(`safePwd2`).String()
		if err := m.CheckNewSafePwd(safePwd); err != nil {
			return safePwd, err
		}
		if safePwd != safePwd2 {
			return safePwd, ctx.E(`两次新密码输入不一致`)
		}
		return safePwd, nil
	default:
		password := ctx.Formx(`password`).String()
		password2 := ctx.Formx(`password2`).String()
		if err := m.CheckNewPassword(password); err != nil {
			return password, err
		}
		if password != password2 {
			return password, ctx.E(`两次新密码输入不一致`)
		}
		return password, nil
	}
}
