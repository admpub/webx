package resetpassword

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/handler/frontend/user/profile"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func mobileValidate(c echo.Context, fieldName string, fieldValue string) error {
	if err := c.Validate(fieldName, fieldValue, `mobile`); err != nil {
		return c.NewError(code.InvalidParameter, `手机号码不正确`).SetZone(fieldName)
	}
	return nil
}

func mobileSend(c echo.Context, m *modelCustomer.Customer, account string) error {
	if len(m.Mobile) == 0 {
		return c.E(`账号"%v"没有设置手机号，不支持通过手机找回密码`, m.Name)
	}
	if m.Mobile != account {
		return c.E(`您输入的手机号与账号设置的不匹配，请输入账号中设置的手机号码`)
	}
	c.Request().Form().Set(`mobile`, account)
	message := `亲爱的客户: {name}，您正在找回密码，点击链接开始重置密码：` + GenResetPasswordURL(m.Name, `mobile`, account) + ` 或者 手动输入验证码为：{code} ({lifeTime}分钟内有效) [{siteName}]`
	return profile.MobileSend(c, m, `forgot`, message)
}

func mobileGetAccount(c echo.Context, m *modelCustomer.Customer) (string, error) {
	if len(m.Mobile) == 0 {
		return m.Mobile, c.E(`账号"%v"没有设置手机号，不支持通过手机找回密码`, m.Name)
	}
	return m.Mobile, nil
}

func mobileOnChangeAfter(c echo.Context, m *modelCustomer.Customer) error {
	if m.MobileBind == `Y` {
		return nil
	}
	set := echo.H{
		`mobile`:      m.Mobile,
		`mobile_bind`: `Y`,
	}
	return m.UpdateFields(nil, set, db.Cond{`id`: m.Id})
}
