package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/webx/application/dbschema"
)

func NewCounter(ctx echo.Context) *Counter {
	m := &Counter{
		OfficialCustomerCounter: dbschema.NewOfficialCustomerCounter(ctx),
	}
	return m
}

type Counter struct {
	*dbschema.OfficialCustomerCounter
}

func (u *Counter) Exists() (bool, error) {
	return u.OfficialCustomerCounter.Exists(nil, db.And(
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`target`: u.Target},
	))
}

func (u *Counter) check() error {
	exists, err := u.Exists()
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `Target已经存在`).SetZone(`target`)
	}
	return err
}

func (u *Counter) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.OfficialCustomerCounter.Insert()
}

func (u *Counter) Incr(target string, n uint64) error {
	exists, err := u.Exists()
	if err != nil {
		return err
	}
	if !exists {
		u.Total = n
		_, err = u.OfficialCustomerCounter.Insert()
		return err
	}
	return u.OfficialCustomerCounter.UpdateField(nil, `total`, db.Raw("total+"+param.AsString(n)), db.And(
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`target`: u.Target},
	))
}

func (u *Counter) Decr(target string, n uint64) error {
	exists, err := u.Exists()
	if err != nil || !exists {
		return err
	}
	return u.OfficialCustomerCounter.UpdateField(nil, `total`, db.Raw("total-"+param.AsString(n)), db.And(
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`target`: u.Target},
	))
}

func (u *Counter) GetTotal(target string) (uint64, error) {
	err := u.OfficialCustomerCounter.Get(nil, db.And(
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`target`: u.Target},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
		return 0, err
	}
	return u.Total, nil
}

func (u *Counter) GetTotals(targets ...string) (map[string]uint64, error) {
	totals := map[string]uint64{}
	if len(targets) == 0 {
		return totals, nil
	}
	for _, v := range targets {
		totals[v] = 0
	}
	_, err := u.OfficialCustomerCounter.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`target`: db.In(targets)},
	))
	if err != nil {
		return totals, err
	}
	for _, row := range u.OfficialCustomerCounter.Objects() {
		totals[row.Target] = row.Total
	}
	return totals, nil
}
