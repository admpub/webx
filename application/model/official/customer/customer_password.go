package customer

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/code"
)

func (f *Customer) CheckSignInPassword(pass string) error {
	if f.Disabled == `Y` {
		return f.Context().NewError(code.UserDisabled, `该用户已被禁用`)
	}
	if f.Password != com.MakePassword(pass, f.Salt) {
		return f.Context().NewError(code.InvalidParameter, `密码不正确`)
	}
	return nil
}

var (
	// CustomerPasswordMinLength 客户登录密码最小长度
	CustomerPasswordMinLength = 8
	// CustomerSafePwdMinLength 客户安全密码最小长度
	CustomerSafePwdMinLength = 6
)

func (f *Customer) CheckNewPassword(pass string) error {
	if len(pass) < CustomerPasswordMinLength {
		return f.Context().NewError(code.InvalidParameter, `新密码不能少于%d个字符`, CustomerPasswordMinLength)
	}
	return nil
}

func (f *Customer) CheckNewSafePwd(pass string) error {
	if len(pass) < CustomerSafePwdMinLength {
		return f.Context().NewError(code.InvalidParameter, `新密码不能少于%d个字符`, CustomerSafePwdMinLength)
	}
	return nil
}

func (f *Customer) UpdateSignInPassword(pass string) error {
	if err := f.CheckNewPassword(pass); err != nil {
		return err
	}
	newPass := com.MakePassword(pass, f.Salt)
	if newPass == f.Password {
		return f.Context().NewError(code.InvalidParameter, `新密码与旧密码一致，无需修改`)
	}
	return f.UpdateField(nil, `password`, com.MakePassword(pass, f.Salt), db.Cond{`id`: f.Id})
}

func (f *Customer) UpdateSafePassword(pass string) error {
	if err := f.CheckNewSafePwd(pass); err != nil {
		return err
	}
	newPass := com.MakePassword(pass, f.Salt)
	if newPass == f.SafePwd {
		return f.Context().NewError(code.InvalidParameter, `新密码与旧密码一致，无需修改`)
	}
	return f.UpdateField(nil, `safe_pwd`, newPass, db.Cond{`id`: f.Id})
}

func (f *Customer) CheckSafePassword(pass string) error {
	if f.SafePwd != com.MakePassword(pass, f.Salt) {
		return f.Context().NewError(code.InvalidParameter, `密码不正确`)
	}
	return nil
}

func (f *Customer) ClearPasswordData(customers ...*dbschema.OfficialCustomer) dbschema.OfficialCustomer {
	if len(customers) > 0 {
		return ClearPasswordData(customers[0])
	}
	return ClearPasswordData(f.OfficialCustomer)
}

func ClearPasswordData(rawCustomer *dbschema.OfficialCustomer) dbschema.OfficialCustomer {
	customer := *(rawCustomer)
	customer.Password = ``
	customer.Salt = ``
	customer.SafePwd = ``
	customer.SessionId = ``
	return customer
}
