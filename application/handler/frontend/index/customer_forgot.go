package index

import (
	"net/url"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/webx/application/handler/frontend/user/profile"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/resetpassword"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

// forgotSendCode 重置密码第一步：发送验证码
func forgotSendCode(c echo.Context) echo.Data {
	data := c.Data()
	m := modelCustomer.NewCustomer(c)
	name := c.Formx(`name`).String()
	if !com.IsUsername(name) {
		return data.SetError(c.NewError(code.UserNotFound, `用户名无效`).SetZone(`name`))
	}
	account := c.Formx(`account`).String()
	typ := c.Formx(`type`).String()
	if len(name) == 0 {
		return data.SetError(c.E(`请输入用户名`))
	}
	if len(typ) == 0 {
		return data.SetError(c.E(`消息类型无效`))
	}
	var target string
	recvType := resetpassword.Get(typ)
	if recvType == nil || !recvType.On {
		return data.SetError(c.E(`不支持的类型: %v`, typ))
	}
	target = recvType.InputName
	err := recvType.Validate(c, `account`, account)
	if err != nil {
		return data.SetError(err)
	}
	if len(account) == 0 {
		return data.SetError(c.E(`请输入%v`, target))
	}
	err = m.Get(nil, `name`, name)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = c.E(`用户不存在: %v`, name)
		}
		return data.SetError(err)
	}
	err = recvType.Send(c, m, account)
	if err != nil {
		return data.SetError(err)
	}
	if common.IsCaptchaErrCode(data.GetCode()) {
		return data
	}
	data.SetInfo(c.T(`验证码已经发送到"%v"，请注意查收`, account))
	return data
}

// forgotSendCode 重置密码第二步：验证验证码并修改密码
func forgotModifyPassword(c echo.Context) echo.Data {
	m := modelCustomer.NewCustomer(c)
	data := c.Data()
	name := c.Formx(`name`).String()
	if !com.IsUsername(name) {
		return data.SetError(c.NewError(code.UserNotFound, `用户名无效`).SetZone(`name`))
	}
	typ := c.Formx(`type`, `email`).String()
	vcode := c.Formx(`vcode`).String()
	if len(name) == 0 {
		return data.SetError(c.E(`请输入用户名`))
	}
	if len(typ) == 0 {
		return data.SetError(c.E(`消息类型无效`))
	}
	if len(vcode) == 0 {
		return data.SetError(c.E(`验证码无效`))
	}
	password := c.Formx(`password`).String()
	if len(password) == 0 {
		return data.SetError(c.E(`请输入新密码`))
	}
	repassword := c.Formx(`repassword`).String()
	if len(repassword) == 0 {
		return data.SetError(c.E(`请确认密码`))
	}
	if password != repassword {
		return data.SetError(c.E(`两次密码输入不一致，请检查`))
	}
	err := m.Get(nil, `name`, name)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = c.E(`用户不存在: %v`, name)
		}
		return data.SetError(err)
	}
	recvType := resetpassword.Get(typ)
	if recvType == nil || !recvType.On {
		return data.SetError(c.E(`不支持的类型: %v`, typ))
	}
	if recvType.OnChangeBefore != nil {
		err = recvType.OnChangeBefore(c, m)
		if err != nil {
			return data.SetError(err)
		}
	}
	var account string
	account, err = recvType.GetAccount(c, m)
	if err != nil {
		return data.SetError(err)
	}
	data = common.VerifyCaptcha(c, frontend.Name, `code`)
	if common.IsFailureCode(data.GetCode()) {
		return data
	}
	vm := model.NewCode(c)
	err = vm.CheckVerificationCode(vcode, `forgot`, m.Id, `customer`, typ, account)
	if err != nil {
		return data.SetError(err)
	}
	err = m.UpdateSignInPassword(password)
	if err != nil {
		return data.SetError(err)
	}
	err = vm.UseVerificationCode(vm.Verification)
	if err != nil {
		return data.SetError(err)
	}
	data.SetInfo(c.T(`密码已经修改成功，请使用新密码登录`))
	if recvType.OnChangeAfter != nil {
		recvType.OnChangeAfter(c, m)
	}
	return data
}

// Forgot 忘记密码
func Forgot(c echo.Context) error {
	var err error
	if c.IsPost() {
		op := c.Form(`op`)
		switch op {
		case `sendCode`: // step 1
			data := forgotSendCode(c)
			return c.JSON(data)
		case `modifyPassword`: // step 2
			data := forgotModifyPassword(c)
			return c.JSON(data)
		default:
			return c.JSON(c.Data().SetError(c.E(`不支持的操作: %v`, op)))
		}
	}

	token := c.Query(`token`)
	if len(token) > 0 {
		// info := `name=` + url.QueryEscape(name) + `&type=` + url.QueryEscape(typ)
		info := common.Crypto().Decode(token)
		if len(info) > 0 {
			urlValues, err := url.ParseQuery(info)
			if err != nil {
				return err
			}
			for key := range urlValues {
				c.Request().Form().Set(key, urlValues.Get(key))
			}
		}
	}
	recvTypes := resetpassword.List()
	typ := c.Form(`type`)
	var recvType *resetpassword.RecvType
	if len(typ) == 0 {
		for _, rType := range recvTypes {
			recvType = rType
			typ = rType.Key
			break
		}
	}
	if recvType == nil {
		recvType = resetpassword.Get(typ)
		if recvType == nil || !recvType.On {
			return c.E(`不支持的类型: %v`, typ)
		}
	}
	tmpl := c.Internal().String(`tmpl`)
	if len(tmpl) == 0 {
		tmpl = `forgot`
	}
	signUpURL := c.Internal().String(`signUpURL`)
	if len(signUpURL) == 0 {
		signUpURL = sessdata.URLFor(`/sign_up`)
	}
	c.Set(`signUpURL`, signUpURL)
	signInURL := c.Internal().String(`signInURL`)
	if len(signInURL) == 0 {
		signInURL = sessdata.URLFor(`/sign_in`)
	}
	c.Set(`signInURL`, signInURL)
	c.Set(`verifyCodeLength`, profile.VerifyCodeLength())
	c.Set(`recvTypes`, recvTypes)
	c.Set(`recvType`, recvType)
	return c.Render(tmpl, handler.Err(c, err))
}
