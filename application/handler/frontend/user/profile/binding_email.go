package profile

import (
	"strings"
	"time"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/cron"
	"github.com/admpub/nging/v4/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/xcommon"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

var (
	verifyCodeLifetime int64 = 5 //验证码有效时长：分钟
	verifyCodeLength   uint  = 8 //验证码长度
)

// VerifyCodeLength 验证码长度
func VerifyCodeLength() uint {
	return verifyCodeLength
}

// 发送验证码邮件
func bindingEmailSend(ctx echo.Context, m *modelCustomer.Customer) error {
	err := EmailSend(ctx, m, `binding`)
	if err != nil {
		return err
	}
	return ctx.JSON(ctx.Data())
}

// EmailSend 发送验证码邮件
func EmailSend(ctx echo.Context, m *modelCustomer.Customer, purpose string, titleAndMessage ...string) error {
	var err error
	now := time.Now()
	vm := model.NewCode(ctx)
	if err := vm.CheckFrequency(
		m.Id,
		`customer`,
		`email`,
		config.Setting().GetStore(`frequency`).GetStore(`email`),
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
	if m.EmailBind != `Y` {
		m.Email = ctx.Formx(`email`).String()
		if !com.IsEmail(m.Email) {
			return ctx.E(`E-mail格式不正确`)
		}
	}

	verifyCode := com.RandomNumeric(VerifyCodeLength())

	toEmail := m.Email
	var toUsername string

	//获取系统配置
	baseCfg := config.Setting().GetStore(`base`)

	//验证码有效期
	lifetime := baseCfg.Int64(`verifyCodeLifetime`, verifyCodeLifetime)
	if lifetime == 0 {
		lifetime = verifyCodeLifetime
	}
	expiry := now.Add(time.Duration(lifetime) * time.Minute)

	siteURL := xcommon.SiteURL(ctx)

	//邮件内容
	title := `[` + baseCfg.String(`siteName`) + `]` + ctx.T(`请查收验证码`)
	content := []byte(ctx.T(`亲爱的客户: %s，您正在进行邮箱验证，本次验证码为：%s (%d分钟内有效)。<br /><br /> 来自：%s<br />时间：%s`, m.Name, verifyCode, lifetime, siteURL+`/`, time.Now().Format(time.RFC3339)))

	if len(titleAndMessage) > 0 {
		placeholders := map[string]string{
			`name`:     m.Name,
			`code`:     verifyCode,
			`lifeTime`: param.AsString(lifetime),
			`siteName`: baseCfg.String(`siteName`),
			`siteURL`:  siteURL + `/`,
			`now`:      time.Now().Format(time.RFC3339),
		}

		switch len(titleAndMessage) {
		case 2:
			message := titleAndMessage[1]
			for find, to := range placeholders {
				message = strings.ReplaceAll(message, `{`+find+`}`, to)
			}
			content = []byte(message)
			fallthrough
		case 1:
			if len(titleAndMessage[0]) > 0 {
				title = titleAndMessage[0]
				for find, to := range placeholders {
					title = strings.ReplaceAll(title, `{`+find+`}`, to)
				}
			}
		}
	}
	vm.Verification.Code = verifyCode
	vm.Verification.OwnerId = m.Id
	vm.Verification.OwnerType = `customer`
	vm.Verification.Purpose = purpose
	vm.Verification.Start = uint(now.Unix())
	vm.Verification.End = uint(expiry.Unix())
	vm.Verification.Disabled = `N`
	vm.Verification.SendMethod = `email`
	vm.Verification.SendTo = toEmail
	if _, addErr := vm.AddVerificationCode(); addErr != nil {
		return addErr
	}
	logM := model.NewSendingLog(ctx)
	logM.Provider = config.Setting().GetStore(`smtp`).String(`engine`)
	logM.Method = `email`
	logM.To = toEmail
	logM.SourceType = `code_verification`
	logM.SourceId = vm.Verification.Id
	logM.Result = ctx.T(`已加入队列`)
	logM.Status = `queued`
	logM.Content = string(content)
	logM.Params = ``
	logM.AppointmentTime = 0
	if _, addErr := logM.Add(); addErr != nil {
		return addErr
	}
	err = cron.SendMailWithID(logM.Id, toEmail, toUsername, title, content)
	if err != nil {
		logM.UpdateFields(nil, echo.H{
			`status`: `failure`,
			`result`: ctx.T(`发送失败`) + `: ` + err.Error(),
		}, `id`, logM.Id)
		return err
	}
	data.SetInfo(`邮件发送成功`)
	return nil
}

// 验证邮件验证码
func bindingEmailVerify(ctx echo.Context, m *modelCustomer.Customer) error {
	data := ctx.Data()
	var operateDesc string
	if m.EmailBind != `Y` {
		m.Email = ctx.Formx(`email`).String()
		if !com.IsEmail(m.Email) {
			return ctx.E(`E-mail格式不正确`)
		}
		m.EmailBind = `Y` //新绑定
		operateDesc = ctx.T(`%s绑定成功`, ctx.T(`邮箱`))
	} else {
		m.EmailBind = `N` //取消原绑定
		operateDesc = ctx.T(`%s解绑成功`, ctx.T(`邮箱`))
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
		return ctx.E(`请输入邮件验证码`)
	}
	vm := model.NewCode(ctx)
	err := vm.CheckVerificationCode(vcode, purpose, m.Id, `customer`, `email`, m.Email)
	if err != nil {
		return err
	}
	return vm.UseVerificationCode(vm.Verification)
}
