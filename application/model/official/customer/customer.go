package customer

import (
	"encoding/gob"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/webx/application/dbschema"
)

func init() {
	gob.Register(dbschema.NewOfficialCustomer(nil))
}

func NewCustomer(ctx echo.Context) *Customer {
	m := &Customer{
		OfficialCustomer: dbschema.NewOfficialCustomer(ctx),
	}
	return m
}

type Customer struct {
	*dbschema.OfficialCustomer
}

func (f *Customer) GetDetail(cond db.Compound) (*CustomerAndGroup, error) {
	m := &CustomerAndGroup{OfficialCustomer: f.OfficialCustomer}
	err := f.Param(func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(m.OfficialCustomer, `password`, `salt`, `safe_pwd`)...).
			Relation(`Level`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
				return sel.Columns(`id`, `short`, `name`, `description`, `icon_image`, `icon_class`, `color`, `bgcolor`)
			}).
			Relation(`AgentLevel`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
				return sel.Columns(`id`, `name`, `description`)
			}).
			Relation(`Roles`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
				return sel.Columns(`id`, `name`, `description`)
			})
	}, cond).SetRecv(m).One()
	return m, err
}

func (f *Customer) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*CustomerAndGroup, error) {
	list := []*CustomerAndGroup{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...).Select(factory.DBIGet().OmitSelect(f.OfficialCustomer, `password`, `salt`, `safe_pwd`)...)
	}, cond.And()).Paging(f.Context())
	return list, err
}

func (f *Customer) Exists(name string, fields ...string) (bool, error) {
	fieldName := `name`
	if len(fields) > 0 {
		fieldName = fields[0]
	}
	cond := db.NewCompounds()
	switch fieldName {
	case `mobile`, `email`:
		cond.Add(db.Cond{fieldName + `_bind`: `Y`})
		fallthrough
	default:
		cond.Add(db.Cond{fieldName: name})
	}
	return f.OfficialCustomer.Exists(nil, cond.And())
}

func (f *Customer) ExistsOther(name string, id uint64, fields ...string) (bool, error) {
	fieldName := `name`
	if len(fields) > 0 {
		fieldName = fields[0]
	}
	cond := db.NewCompounds()
	switch fieldName {
	case `mobile`, `email`:
		cond.Add(db.Cond{fieldName + `_bind`: `Y`})
		fallthrough
	default:
		cond.Add(db.Cond{fieldName: name})
		cond.Add(db.Cond{`id <>`: id})
	}
	return f.OfficialCustomer.Exists(nil, cond.And())
}

func (f *Customer) check(isNew bool, id uint64) (err error) {
	if len(f.Name) == 0 {
		err = f.Context().NewError(code.InvalidParameter, `请输入登录名称`)
		return
	}
	if !com.IsUsername(f.Name) {
		return f.Context().NewError(code.InvalidParameter, `用户名格式不正确(只能包含字母、数字、下划线和汉字)`)
	}
	var exists bool
	if len(f.Mobile) > 0 {
		if _err := f.Context().Validate(`mobile`, f.Mobile, `mobile`); _err != nil {
			err = f.Context().NewError(code.InvalidParameter, `手机号码格式不正确`)
			return
		}
		if isNew {
			exists, err = f.Exists(f.Mobile, `mobile`)
		} else {
			exists, err = f.ExistsOther(f.Mobile, id, `mobile`)
		}
		if err != nil {
			return
		}
		if exists {
			err = f.Context().NewError(code.DataAlreadyExists, `手机号码“%s”已经存在`, f.Mobile)
			return
		}
	}
	if len(f.Email) > 0 {
		if _err := f.Context().Validate(`email`, f.Email, `email`); _err != nil {
			err = f.Context().NewError(code.InvalidParameter, `E-mail格式不正确`)
			return
		}
		if isNew {
			exists, err = f.Exists(f.Email, `email`)
		} else {
			exists, err = f.ExistsOther(f.Mobile, id, `email`)
		}
		if err != nil {
			return
		}
		if exists {
			err = f.Context().NewError(code.DataAlreadyExists, `E-mail地址“%s”已经存在`, f.Email)
			return
		}
	}

	if isNew {
		exists, err = f.Exists(f.Name)
	} else {
		exists, err = f.ExistsOther(f.Mobile, id)
	}
	if err != nil {
		return
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `用户名称“%s”已经存在`, f.Name)
		return
	}
	f.MobileBind = common.GetBoolFlag(f.MobileBind)
	f.EmailBind = common.GetBoolFlag(f.EmailBind)
	f.Disabled = common.GetBoolFlag(f.Disabled)
	f.Online = common.GetBoolFlag(f.Online)
	return nil
}

func (f *Customer) Add() (pk interface{}, err error) {
	err = f.check(true, 0)
	if err != nil {
		return
	}
	if len(f.Password) == 0 {
		err = f.Context().NewError(code.InvalidParameter, `请输入登录密码`)
		return
	}
	f.Salt = com.Salt()
	f.Password = com.MakePassword(f.Password, f.Salt)
	if len(f.SafePwd) > 0 {
		f.SafePwd = com.MakePassword(f.SafePwd, f.Salt)
	}
	pk, err = f.OfficialCustomer.Insert()
	return
}

func (f *Customer) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	m := dbschema.NewOfficialCustomer(nil)
	m.CPAFrom(f.OfficialCustomer)
	err := m.Get(nil, args...)
	if err != nil {
		return err
	}
	f.Name = m.Name
	f.Salt = m.Salt
	if len(f.Password) == 0 {
		f.Password = m.Password
	} else if f.Password != m.Password {
		f.Password = com.MakePassword(f.Password, m.Salt)
	}
	if len(f.SafePwd) == 0 {
		f.SafePwd = m.SafePwd
	} else if f.SafePwd != m.SafePwd {
		f.SafePwd = com.MakePassword(f.SafePwd, m.Salt)
	}
	err = f.check(false, m.Id)
	if err != nil {
		return err
	}
	err = f.OfficialCustomer.Update(mw, args...)
	return err
}

func (f *Customer) NewLoginLog(username string, authType string) *model.LoginLog {
	loginLogM := model.NewLoginLog(f.Context())
	loginLogM.OwnerType = `customer`
	loginLogM.Username = username
	loginLogM.Success = `N`
	loginLogM.SessionId = f.Context().Session().MustID()
	loginLogM.AuthType = authType
	return loginLogM
}

// APIData 返回给API接口的数据
func (f *Customer) APIData(customers ...*dbschema.OfficialCustomer) echo.H {
	customer := f.OfficialCustomer
	if len(customers) > 0 {
		customer = customers[0]
	}
	claims := echo.H{}
	claims[`id`] = customer.Id
	claims[`name`] = customer.Name
	claims[`name`] = customer.Name
	claims[`avatar`] = customer.Avatar
	claims[`gender`] = customer.Gender
	claims[`groupId`] = customer.GroupId
	claims[`realName`] = customer.RealName
	claims[`mobile`] = customer.Mobile
	claims[`mobileBind`] = customer.MobileBind
	claims[`email`] = customer.Email
	claims[`emailBind`] = customer.EmailBind
	claims[`levelId`] = customer.LevelId
	claims[`agentLevel`] = customer.AgentLevel
	claims[`inviterId`] = customer.InviterId
	return claims
}
