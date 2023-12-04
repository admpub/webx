package customer

import (
	"strconv"

	"github.com/admpub/goth"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/oauth2client"
)

func NewOAuth(ctx echo.Context) *OAuth {
	m := &OAuth{
		OfficialCustomerOauth: dbschema.NewOfficialCustomerOauth(ctx),
	}
	return m
}

type OAuth struct {
	*dbschema.OfficialCustomerOauth
}

func (f *OAuth) Add() (pk interface{}, err error) {
	old := dbschema.NewOfficialCustomerOauth(f.Context())
	err = old.Get(nil, db.And(
		db.Cond{`customer_id`: f.CustomerId},
		db.Cond{`union_id`: f.UnionId},
		db.Cond{`open_id`: f.OpenId},
		db.Cond{`type`: f.Type},
	))
	if err == nil {
		pk = old.Id
		set := echo.H{}
		if len(f.Email) > 0 && old.Email != f.Email {
			set[`email`] = f.Email
		}
		if len(f.Mobile) > 0 && old.Mobile != f.Mobile {
			set[`mobile`] = f.Mobile
		}
		if len(f.Avatar) > 0 && old.Avatar != f.Avatar {
			set[`avatar`] = f.Avatar
		}
		if len(f.AccessToken) > 0 && old.AccessToken != f.AccessToken {
			set[`access_token`] = f.AccessToken
		}
		if len(f.RefreshToken) > 0 && old.RefreshToken != f.RefreshToken {
			set[`refresh_token`] = f.RefreshToken
		}
		if f.Expired > 0 && old.Expired != f.Expired {
			set[`expired`] = f.Expired
		}
		if len(set) == 0 {
			return
		}
		err = f.UpdateFields(nil, set, `id`, old.Id)
		return
	}
	if err != db.ErrNoMoreRows {
		return
	}
	return f.OfficialCustomerOauth.Insert()
}

func (f *OAuth) Upsert(mw func(db.Result) db.Result, args ...interface{}) (interface{}, error) {
	return f.OfficialCustomerOauth.Upsert(mw, args...)
}

func (f *OAuth) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCustomerOauth.Update(mw, args...)
}

func (f *OAuth) GetByOutUser(user *goth.User) (err error) {
	var unionID string
	if v, y := user.RawData[`unionid`]; y {
		unionID, _ = v.(string)
	}
	return f.Get(nil, db.And(
		db.Cond{`union_id`: unionID},
		db.Cond{`open_id`: user.UserID},
		db.Cond{`type`: user.Provider},
	))
}

func (f *OAuth) CopyFrom(user *goth.User) *OAuth {
	var unionID string
	if v, y := user.RawData[`unionid`]; y {
		unionID, _ = v.(string)
	}
	f.UnionId = unionID
	f.OpenId = user.UserID
	f.Type = user.Provider
	f.AccessToken = user.AccessToken
	f.RefreshToken = user.RefreshToken
	f.Expired = uint(user.ExpiresAt.Unix())
	return f
}

func (f *OAuth) Exists(customerID uint64, unionID string, openID string, typ string) (bool, error) {
	return f.OfficialCustomerOauth.Exists(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`union_id`: unionID},
		db.Cond{`open_id`: openID},
		db.Cond{`type`: typ},
	))
}

func (f *OAuth) ExistsOtherBinding(customerID uint64, id uint64) (bool, error) {
	return f.OfficialCustomerOauth.Exists(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`id`: db.NotEq(id)},
	))
}

func (f *OAuth) OAuthUserGender(ouser *goth.User) string {
	if v, y := ouser.RawData[`gender`]; y {
		gender := com.String(v)
		if len(gender) > 0 {
			switch gender[0] {
			case 'F', 'f', '0':
				return `female`
			case 'M', 'm', '1':
				return `male`
			default:
				return `secret` //保密
			}
		}
	}
	return ``
}

// SignUpCustomer 用户不存在，需要新注册
func (f *OAuth) SignUpCustomer(ouser *goth.User) (*Customer, error) {
	customerM := NewCustomer(f.Context())
	customerM.Name = ouser.Provider + `_` + com.RandomAlphanumeric(5)
	customerM.Password = com.RandomString(16)
	customerM.Email = ouser.Email
	customerM.Avatar = ouser.AvatarURL
	customerM.RegisteredBy = `oauth2.login`

	username := customerM.Name
	exists, err := customerM.Exists(customerM.Name)
	if err != nil {
		return customerM, err
	}
	for i := 0; exists; i++ {
		customerM.Name = username + com.RandomAlphanumeric(1) + strconv.Itoa(i)
		exists, err = customerM.Exists(customerM.Name)
		if err != nil {
			return customerM, err
		}
	}
	customerM.Gender = f.OAuthUserGender(ouser)
	//注册并登录
	err = customerM.SignUp(customerM.Name, customerM.Password, customerM.Mobile, customerM.Email)
	if err != nil {
		return customerM, err
	}
	return customerM, err
}

func (f *OAuth) SaveSession(ouser *goth.User) error {
	return oauth2client.SaveSession(f.Context(), ouser)
}

func (f *OAuth) GetSession() (*goth.User, bool, error) {
	return oauth2client.GetSession(f.Context())
}

func (f *OAuth) DelSession() {
	oauth2client.DelSession(f.Context())
}
