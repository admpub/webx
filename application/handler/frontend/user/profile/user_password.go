package profile

import (
	"errors"

	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/library/sendmsg"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

type StepID int

func (s StepID) String() string {
	k := int(s)
	if len(PasswordModifySteps) <= k {
		return ``
	}
	return PasswordModifySteps[k]
}

func (s StepID) Int() int {
	return int(s)
}

const (
	StepInit StepID = iota
	StepSend
	StepVerifyAndModify
)

var (
	PasswordModifyStepIDs = map[string]StepID{
		`init`:   StepInit,
		`send`:   StepSend,
		`verify`: StepVerifyAndModify,
	}
	PasswordModifySteps = []string{
		`init`,
		`send`,
		`verify`,
	}
)

// Password 修改密码
func Password(c echo.Context) error {
	typ := c.Form(`type`, `sign_in`)
	typName := c.T(`登录密码`)
	if typ != `sign_in` {
		typName = c.T(`安全密码`)
	}
	customer := xMW.Customer(c)
	m := modelCustomer.NewCustomer(c)
	err := m.VerifySession(customer)
	if err != nil {
		if common.IsUserNotLoggedIn(err) {
			return c.E(`请先登录`)
		}
		return err
	}
	if c.IsPost() {
		var newPassword string
		newPassword, err = checkNewPassword(c, m, typ)
		if err != nil {
			return err
		}
		authType := `password`
		if m.MobileBind == `Y` || m.EmailBind == `Y` {
			atype := c.Form(`authType`)
			switch atype {
			case `email`, `mobile`:
				authType = atype
			default:
				if m.MobileBind == `Y` {
					authType = `mobile`
				} else {
					authType = `email`
				}
			}
		}
		step := c.Form(`step`, `init`)
		verify := c.Formx(`verify`).Bool()
		if verify {
			step = `verify`
		}
		stepID, ok := PasswordModifyStepIDs[step]
		if !ok {
			stepID = StepInit
			step = stepID.String()
		}
		data := c.Data()
		genRespData := func() echo.H {
			stepID++
			return echo.H{
				`authType`: authType,
				`nextStep`: stepID.String(),
			}
		}
		switch authType {
		case `password`:
			if stepID == StepVerifyAndModify {
				oldPwd := c.Formx(`oldPwd`).String()
				if len(oldPwd) == 0 {
					return errors.New(`请输入旧密码`)
				}
				err = m.CheckSignInPassword(oldPwd)
				if err != nil {
					return err
				}
				err = modifyPassword(m, typ, newPassword)
				if err != nil {
					return err
				}
				data.SetInfo(c.T(`您的%s已经修改成功`, typName))
			}
		case `email`:
			if stepID == StepSend {
				err = sendmsg.EmailSend(c, m, `modify-password`)
				if err != nil {
					return err
				}
				data.GetData().(echo.H).DeepMerge(genRespData())
				return c.JSON(data)
			}
			if stepID == StepVerifyAndModify {
				err = EmailVerify(c, m, `modify-password`)
				if err != nil {
					return err
				}
				err = modifyPassword(m, typ, newPassword)
				if err != nil {
					return err
				}
				data.SetInfo(c.T(`您的%s已经修改成功`, typName))
			}
		case `mobile`:
			if stepID == StepSend {
				err = sendmsg.MobileSend(c, m, `modify-password`)
				if err != nil {
					return err
				}
				data.GetData().(echo.H).DeepMerge(genRespData())
				return c.JSON(data)
			}
			if stepID == StepVerifyAndModify {
				err = MobileVerify(c, m, `modify-password`)
				if err != nil {
					return err
				}
				err = modifyPassword(m, typ, newPassword)
				if err != nil {
					return err
				}
				data.SetInfo(c.T(`您的%s已经修改成功`, typName))
			}
		}
		data.SetData(genRespData())
		return c.JSON(data)
	}
	c.Set(`activeURL`, `/user/profile`)
	c.Set(`type`, typ)
	c.Set(`passwordMinLength`, modelCustomer.CustomerPasswordMinLength)
	c.Set(`safePwdMinLength`, modelCustomer.CustomerSafePwdMinLength)
	c.Set(`title`, c.T(`修改密码`))
	return c.Render(`/user/profile/password`, err)
}
