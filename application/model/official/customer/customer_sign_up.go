package customer

import (
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/model"
	modelLevel "github.com/admpub/webx/application/model/official/level"
)

// SignUp 注册用户
func (f *Customer) SignUp(user, pass, mobile, email string, options ...CustomerOption) error {
	co := NewCustomerOptions(f.OfficialCustomer, true)
	co.Name = user
	co.Password = pass
	co.Mobile = mobile
	co.Email = email
	for _, option := range options {
		option(co)
	}
	if f.LevelId < 1 {
		levelM := modelLevel.NewLevel(f.Context())
		if level, err := levelM.CanAutoLevelUpByIntegral(0); err == nil {
			f.LevelId = level.Id
		}
	}
	if len(f.RoleIds) == 0 {
		roleM := NewRole(f.Context())
		if err := roleM.GetDefault(); err == nil {
			f.RoleIds = param.AsString(roleM.Id)
		}
	}
	f.SessionId = f.Context().Session().MustID()
	_, err := f.Add()
	if err != nil {
		return err
	}

	return f.FireSignUpSuccess(co, model.AuthTypePassword, options...)
}

func (f *Customer) FireSignUpSuccess(co *CustomerOptions, authType string, options ...CustomerOption) (err error) {
	integral := config.Setting(`base`, `addExperience`).Float64(`register`)
	if err = f.AddRewardOnSignUp(integral); err != nil {
		return err
	}
	err = f.LinkOAuthUser()
	if err != nil {
		return err
	}
	err = FireSignIn(f.OfficialCustomer)
	if err != nil {
		return err
	}
	deviceM := NewDevice(f.Context())
	deviceM.SessionId = f.SessionId
	deviceM.CustomerId = f.Id
	deviceM.SetOptions(co)
	_, err = deviceM.Upsert()
	if err != nil {
		return err
	}

	f.SetSession()

	loginLogM := f.NewLoginLog(co.Name, authType)
	loginLogM.Success = `Y`
	loginLogM.AddAndSaveSession()
	return
}
