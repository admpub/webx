package profile

import (
	"fmt"
	"strings"
	"time"

	"github.com/coscms/sms"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/model"
	uploadChecker "github.com/admpub/nging/v5/application/registry/upload/checker"
	"github.com/admpub/webx/application/initialize/frontend"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

var _ = xMW.URLFor

// 验证短信验证码
func bindingMobileVerify(ctx echo.Context, m *modelCustomer.Customer) error {
	data := ctx.Data()
	var operateDesc string
	if m.MobileBind != `Y` {
		m.Mobile = ctx.Formx(`mobile`).String()
		if !ctx.Validate(`mobile`, m.Mobile, `mobile`).Ok() {
			return ctx.E(`手机号码格式不正确`)
		}
		m.MobileBind = `Y` //新绑定
		operateDesc = ctx.T(`%s绑定成功`, ctx.T(`手机`))
	} else {
		m.MobileBind = `N` //取消原绑定
		operateDesc = ctx.T(`%s解绑成功`, ctx.T(`手机`))
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
		return ctx.E(`请输入短信验证码`)
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
	err := MobileSend(ctx, m, `binding`)
	if err != nil {
		return err
	}
	return ctx.JSON(ctx.Data())
}

// MobileSend 发送验证码短信
func MobileSend(ctx echo.Context, m *modelCustomer.Customer, purpose string, messages ...string) error {
	var err error
	vm := model.NewCode(ctx)
	now := time.Now()

	if err := vm.CheckFrequency(
		m.Id,
		`customer`,
		`mobile`,
		config.Setting().GetStore(`frequency`).GetStore(`mobile`),
	); err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
	}
	vm.Verification.Reset()
	data := common.VerifyCaptcha(ctx, frontend.Name, `code`)
	if common.IsFailureCode(data.GetCode()) {
		return nil
	}
	if m.MobileBind != `Y` {
		m.Mobile = ctx.Formx(`mobile`).String()
		if !ctx.Validate(`mobile`, m.Mobile, `mobile`).Ok() {
			return ctx.E(`手机号码格式不正确`)
		}
	}

	verifyCode := com.RandomNumeric(VerifyCodeLength())
	//发送短信
	provider, smsProviderName := sms.AnyOne()
	if provider == nil || len(smsProviderName) == 0 {
		err = ctx.E(`找不到短信发送服务`)
		return err
	}
	smsConfig := sms.NewConfig()
	smsConfig.Mobile = m.Mobile

	//获取系统配置
	baseCfg := config.Setting().GetStore(`base`)

	//验证码有效期
	lifetime := baseCfg.Int64(`verifyCodeLifetime`, verifyCodeLifetime)
	if lifetime == 0 {
		lifetime = verifyCodeLifetime
	}
	expiry := now.Add(time.Duration(lifetime) * time.Minute)
	var message string
	if len(messages) > 0 {
		message = messages[0]
	}
	if len(message) == 0 {
		message = ctx.T(`亲爱的客户: %s，您正在进行手机号码验证，本次验证码为：%s (%d分钟内有效) [%s]`, m.Name, verifyCode, lifetime, baseCfg.String(`siteName`))
	} else {
		placeholders := map[string]string{
			`name`:     m.Name,
			`code`:     verifyCode,
			`lifeTime`: param.AsString(lifetime),
			`siteName`: baseCfg.String(`siteName`),
		}
		for find, to := range placeholders {
			message = strings.ReplaceAll(message, `{`+find+`}`, to)
		}
	}
	smsConfig.Content = message
	smsConfig.Template = ``
	smsConfig.SignName = ``
	smsConfig.ExtraData[`code`] = verifyCode

	// 记录日志
	vm.Verification.Code = verifyCode
	vm.Verification.OwnerId = m.Id
	vm.Verification.OwnerType = `customer`
	vm.Verification.Purpose = purpose
	vm.Verification.Start = uint(now.Unix())
	vm.Verification.End = uint(expiry.Unix())
	vm.Verification.Disabled = `N`
	vm.Verification.SendMethod = `mobile`
	vm.Verification.SendTo = m.Mobile
	if _, addErr := vm.AddVerificationCode(); addErr != nil {
		return addErr
	}
	logM := model.NewSendingLog(ctx)
	logM.Provider = smsProviderName
	logM.Method = `mobile`
	logM.To = m.Mobile
	logM.SourceType = `code_verification`
	logM.SourceId = vm.Verification.Id
	logM.Result = ctx.T(`发送成功`)
	logM.Status = `success`
	logM.Content = smsConfig.Content
	b, e := com.JSONEncode(smsConfig.ExtraData)
	if e != nil {
		return e
	}
	logM.Params = string(b)
	logM.AppointmentTime = 0
	if _, addErr := logM.Add(); addErr != nil {
		return addErr
	}
	timestamp := time.Now().Unix()
	smsConfig.CallbackURL = xMW.URLFor(`/verification/callback/` + smsProviderName + `/` + fmt.Sprint(vm.Verification.Id) + `/` + fmt.Sprint(timestamp) + `/` + uploadChecker.Token(smsProviderName, vm.Verification.Id, timestamp))

	err = provider.Send(smsConfig)
	if err != nil {
		logM.UpdateFields(nil, echo.H{
			`status`: `failure`,
			`result`: ctx.T(`发送失败`) + `: ` + err.Error(),
		}, `id`, logM.Id)
		return err
	}
	data.SetInfo(ctx.T(`短信发送成功`))
	return nil
}
