package profile

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/model"
	"github.com/coscms/webfront/library/sendmsg"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

var _ = xMW.URLFor

// 验证短信验证码
func bindingMobileVerify(ctx echo.Context, m *modelCustomer.Customer) error {
	data := ctx.Data()
	var operateDesc string
	if m.MobileBind != `Y` {
		m.Mobile = ctx.Formx(`mobile`).String()
		if err := ctx.Validate(`mobile`, m.Mobile, `mobile`); err != nil {
			return ctx.NewError(code.InvalidParameter, `手机号码格式不正确`).SetZone(`mobile`)
		}
		customerM := modelCustomer.NewCustomer(ctx)
		exists, err := customerM.ExistsOther(m.Mobile, m.Id, `mobile`)
		if err != nil {
			return err
		}
		if exists {
			return ctx.NewError(code.DataAlreadyExists, `手机号码“%s”已被其他账号绑定`, m.Mobile).SetZone(`mobile`)
		}
		m.MobileBind = `Y` //新绑定
		operateDesc = ctx.T(`%s绑定成功`, ctx.T(`手机号码`))
	} else {
		m.MobileBind = `N` //取消原绑定
		operateDesc = ctx.T(`%s解绑成功`, ctx.T(`手机号码`))
	}
	ctx.Begin()
	err := MobileVerify(ctx, m, `binding`)
	if err != nil {
		ctx.Rollback()
		return err
	}
	set := echo.H{
		`mobile`:      m.Mobile,
		`mobile_bind`: m.MobileBind,
	}
	err = m.UpdateFields(nil, set, db.Cond{`id`: m.Id})
	if err != nil {
		ctx.Rollback()
		return err
	}
	ctx.Commit()
	m.SetSession()
	return ctx.JSON(data.SetInfo(operateDesc))
}

// MobileVerify 验证短信
func MobileVerify(ctx echo.Context, m *modelCustomer.Customer, purpose string) error {
	vcode := ctx.Formx(`vcode`).String()
	if len(vcode) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入短信验证码`).SetZone(`vcode`)
	}
	vm := model.NewCode(ctx)
	err := vm.CheckVerificationCode(vcode, purpose, m.Id, `customer`, `mobile`, m.Mobile)
	if err != nil {
		return err
	}
	return vm.UseVerificationCode(vm.Verification)
}

// 发送验证码短信
func bindingMobileSend(ctx echo.Context, m *modelCustomer.Customer) error {
	err := sendmsg.MobileSend(ctx, m, `binding`)
	if err != nil {
		return err
	}
	return ctx.JSON(ctx.Data())
}
