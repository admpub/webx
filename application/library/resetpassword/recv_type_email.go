package resetpassword

import (
	"github.com/admpub/webx/application/handler/frontend/user/profile"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func emailValidate(c echo.Context, fieldName string, fieldValue string) error {
	if err := c.Validate(fieldName, fieldValue, `email`); err != nil {
		return c.NewError(code.InvalidParameter, `E-mail地址不正确`).SetZone(fieldName)
	}
	return nil
}

func emailSend(c echo.Context, m *modelCustomer.Customer, account string) error {
	if len(m.Email) == 0 {
		return c.E(`账号"%v"没有设置E-mail地址，不支持通过E-mail找回密码`, m.Name)
	}
	if m.Email != account {
		return c.E(`您输入的E-mail地址与账号设置的不匹配，请输入账号中设置的E-mail地址`)
	}
	title := `[{siteName}]` + c.T(`找回密码`)
	resetPasswordURL := GenResetPasswordURL(m.Name, `email`, account)
	c.Request().Form().Set(`email`, account)
	message := `亲爱的客户: {name}，您正在找回密码，点击链接开始重置密码：<a href="` + resetPasswordURL + `" target="_blank">` + resetPasswordURL + `</a> ({lifeTime}分钟内有效)。<br /><br /> 来自：{siteURL}<br />时间：{now}`
	return profile.EmailSend(c, m, `forgot`, title, message)
}

func emailGetAccount(c echo.Context, m *modelCustomer.Customer) (string, error) {
	if len(m.Email) == 0 {
		return m.Email, c.E(`账号"%v"没有设置E-mail地址，不支持通过E-mail找回密码`, m.Name)
	}
	return m.Email, nil
}

func emailOnChangeAfter(c echo.Context, m *modelCustomer.Customer) error {
	if m.EmailBind == `Y` {
		return nil
	}
	set := echo.H{
		`email`:      m.Email,
		`email_bind`: `Y`,
	}
	return m.UpdateFields(nil, set, db.Cond{`id`: m.Id})
}
