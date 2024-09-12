package profile

import (
	"github.com/coscms/webcore/model"
	"github.com/coscms/webfront/library/sendmsg"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

// 发送验证码邮件
func bindingEmailSend(ctx echo.Context, m *modelCustomer.Customer) error {
	err := sendmsg.EmailSend(ctx, m, `binding`)
	if err != nil {
		return err
	}
	return ctx.JSON(ctx.Data())
}

// 验证邮件验证码
func bindingEmailVerify(ctx echo.Context, m *modelCustomer.Customer) error {
	data := ctx.Data()
	var operateDesc string
	if m.EmailBind != `Y` {
		m.Email = ctx.Formx(`email`).String()
		if !com.IsEmail(m.Email) {
			return ctx.NewError(code.InvalidParameter, `E-mail格式不正确`).SetZone(`email`)
		}
		customerM := modelCustomer.NewCustomer(ctx)
		exists, err := customerM.ExistsOther(m.Email, m.Id, `email`)
		if err != nil {
			return err
		}
		if exists {
			return ctx.NewError(code.DataAlreadyExists, `E-mail地址“%s”已被其他账号绑定`, m.Email).SetZone(`email`)
		}
		m.EmailBind = `Y` //新绑定
		operateDesc = ctx.T(`%s绑定成功`, ctx.T(`邮箱地址`))
	} else {
		m.EmailBind = `N` //取消原绑定
		operateDesc = ctx.T(`%s解绑成功`, ctx.T(`邮箱地址`))
	}
	ctx.Begin()
	err := EmailVerify(ctx, m, `binding`)
	if err != nil {
		ctx.Rollback()
		return err
	}
	set := echo.H{
		`email`:      m.Email,
		`email_bind`: m.EmailBind,
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

// EmailVerify 验证邮件
func EmailVerify(ctx echo.Context, m *modelCustomer.Customer, purpose string) error {
	vcode := ctx.Formx(`vcode`).String()
	if len(vcode) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入邮件验证码`).SetZone(`vcode`)
	}
	vm := model.NewCode(ctx)
	err := vm.CheckVerificationCode(vcode, purpose, m.Id, `customer`, `email`, m.Email)
	if err != nil {
		return err
	}
	return vm.UseVerificationCode(vm.Verification)
}
