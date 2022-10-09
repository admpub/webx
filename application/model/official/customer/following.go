package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func NewFollowing(ctx echo.Context) *Following {
	m := &Following{
		OfficialCustomerFollowing: dbschema.NewOfficialCustomerFollowing(ctx),
	}
	return m
}

type Following struct {
	*dbschema.OfficialCustomerFollowing
}

func (f *Following) Add() (pk interface{}, err error) {
	var exists bool
	f.Context().Begin()
	exists, err = f.Exists(f.CustomerA, f.CustomerB)
	if err != nil {
		f.Context().Rollback()
		return
	}
	if exists {
		err = f.Context().E(`您已经关注过了`)
		f.Context().Rollback()
		return
	}
	exists, err = f.Exists(f.CustomerB, f.CustomerA)
	if err != nil {
		f.Context().Rollback()
		return
	}
	if exists {
		f.Mutual = `Y`
		err = f.UpdateField(nil, `mutual`, `Y`, db.And(
			db.Cond{`customer_a`: f.CustomerB},
			db.Cond{`customer_b`: f.CustomerA},
		))
		if err != nil {
			f.Context().Rollback()
			return
		}
	} else {
		f.Mutual = `N`
	}
	pk, err = f.OfficialCustomerFollowing.Insert()
	if err != nil {
		f.Context().Rollback()
		return
	}
	customerM := dbschema.OfficialCustomer{}
	err = customerM.UpdateField(nil, `following`, db.Raw(`following+1`), db.Cond{`id`: f.CustomerA})
	if err != nil {
		f.Context().Rollback()
		return
	}
	err = customerM.UpdateField(nil, `followers`, db.Raw(`followers+1`), db.Cond{`id`: f.CustomerB})
	if err != nil {
		f.Context().Rollback()
		return
	}
	f.Context().End(true)
	return
}

func (f *Following) Exists(customerA uint64, customerB uint64) (bool, error) {
	return f.OfficialCustomerFollowing.Exists(nil, db.And(
		db.Cond{`customer_a`: customerA},
		db.Cond{`customer_b`: customerB},
	))
}

func (f *Following) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCustomerFollowing.Update(mw, args...)
}

func (f *Following) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	f.Context().Begin()
	err := f.Get(nil, args...)
	if err != nil {
		return err
	}
	exists, err := f.Exists(f.CustomerB, f.CustomerA)
	if err != nil {
		f.Context().Rollback()
		return err
	}
	if exists {
		err = f.UpdateField(nil, `mutual`, `N`, db.And(
			db.Cond{`customer_a`: f.CustomerB},
			db.Cond{`customer_b`: f.CustomerA},
		))
		if err != nil {
			f.Context().Rollback()
			return err
		}
	}
	err = f.OfficialCustomerFollowing.Delete(mw, args...)
	if err != nil {
		f.Context().Rollback()
		return err
	}

	customerM := dbschema.OfficialCustomer{}
	err = customerM.UpdateField(nil, `following`, db.Raw(`following-1`), db.Cond{`id`: f.CustomerA})
	if err != nil {
		f.Context().Rollback()
		return err
	}
	err = customerM.UpdateField(nil, `followers`, db.Raw(`followers-1`), db.Cond{`id`: f.CustomerB})
	if err != nil {
		f.Context().Rollback()
		return err
	}
	f.Context().End(err == nil)
	return err
}
