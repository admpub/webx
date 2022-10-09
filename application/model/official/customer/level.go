package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

// NewLevel 客户关联等级信息
func NewLevel(ctx echo.Context) *Level {
	m := &Level{
		OfficialCustomerLevelRelation: dbschema.NewOfficialCustomerLevelRelation(ctx),
	}
	return m
}

type Level struct {
	*dbschema.OfficialCustomerLevelRelation
}

func (f *Level) HasLevel(customerID uint64, levelIds ...interface{}) (bool, error) {
	err := f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`status`: `success`},
		db.Cond{`level_id`: db.In(levelIds)},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (f *Level) Add() (pk interface{}, err error) {
	m := dbschema.NewOfficialCustomerLevelRelation(f.Context())
	err = m.Get(nil, db.And(
		db.Cond{`customer_id`: f.CustomerId},
		db.Cond{`level_id`: f.LevelId},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		pk, err = f.OfficialCustomerLevelRelation.Insert()
		return
	}
	err = f.Edit(nil, `id`, m.Id)
	return
}

func (f *Level) Exists(customerID uint64, levelID uint) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV(`customer_id`, customerID)
	cond.AddKV(`level_id`, levelID)
	return f.OfficialCustomerLevelRelation.Exists(nil, cond.And())
}

func (f *Level) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCustomerLevelRelation.Update(mw, args...)
}

func (f *Level) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCustomerLevelRelation.Delete(mw, args...)
	return err
}
