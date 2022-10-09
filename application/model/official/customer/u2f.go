package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
)

func NewU2F(ctx echo.Context) *U2F {
	m := &U2F{
		OfficialCustomerU2f: dbschema.NewOfficialCustomerU2f(ctx),
	}
	return m
}

type U2F struct {
	*dbschema.OfficialCustomerU2f
}

func (u *U2F) check() error {
	exists, err := u.Exists(nil, db.And(
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`type`: u.Type},
		db.Cond{`step`: u.Step},
		db.Cond{`token`: u.Token},
	))
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `Token已经存在`).SetZone(`token`)
	}
	return err
}

func (u *U2F) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.OfficialCustomerU2f.Insert()
}

func (u *U2F) HasType(customerID uint64, authType string, step uint) (bool, error) {
	return u.OfficialCustomerU2f.Exists(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`type`: authType},
		db.Cond{`step`: step},
	))
}

func (u *U2F) Unbind(customerID uint64, authType string, step uint) error {
	return u.OfficialCustomerU2f.Delete(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`type`: authType},
		db.Cond{`step`: step},
	))
}

func (u *U2F) UnbindByToken(customerID uint64, authType string, step uint, token string) error {
	return u.OfficialCustomerU2f.Delete(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`type`: authType},
		db.Cond{`step`: step},
		db.Cond{`token`: token},
	))
}

func (u *U2F) ListPageByType(customerID uint64, authType string, step uint, sorts ...interface{}) error {
	cond := db.NewCompounds()
	cond.AddKV(`customer_id`, customerID)
	cond.AddKV(`type`, authType)
	cond.AddKV(`step`, step)
	return u.ListPage(cond, sorts...)
}
