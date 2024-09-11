package customer

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/common"
	"github.com/admpub/webx/application/dbschema"
)

func NewGroupPackage(ctx echo.Context) *GroupPackage {
	m := &GroupPackage{
		OfficialCustomerGroupPackage: dbschema.NewOfficialCustomerGroupPackage(ctx),
	}
	return m
}

// GroupPackage 等级组套餐价格
type GroupPackage struct {
	*dbschema.OfficialCustomerGroupPackage
}

func (u *GroupPackage) check() error {
	exists, err := u.Exists(nil, db.And(
		db.Cond{`group`: u.Group},
		db.Cond{`price`: u.Price},
		db.Cond{`time_duration`: u.TimeDuration},
		db.Cond{`time_unit`: u.TimeUnit},
	))
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `相同时段、价格的套餐已经存在`)
	}
	return err
}

func (u *GroupPackage) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.OfficialCustomerGroupPackage.Insert()
}

func (f *GroupPackage) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCustomerGroupPackage.Update(mw, args...)
}

func (f *GroupPackage) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCustomerGroupPackage.Delete(mw, args...)
	return err
}

// ListByGroup 列出某个组的套餐信息
func (u *GroupPackage) ListByGroup(group string) error {
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, common.BoolN)
	cond.AddKV(`group`, group)
	_, err := u.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, -1, cond.And())
	return err
}

func (u *GroupPackage) IncrSold(id uint, n ...int) error {
	var _n int
	if len(n) > 0 {
		_n = n[0]
	}
	if _n == 0 {
		_n = 1
	}
	return u.UpdateField(nil, `sold`, db.Raw(`sold+`+param.AsString(_n)), `id`, id)
}

// ListGroup 列出设置有套餐的组名称
func (u *GroupPackage) ListGroup() ([]string, error) {
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, common.BoolN)
	_, err := u.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`group`).Group(`group`)
	}, 0, -1, cond.And())
	if err != nil {
		return nil, err
	}
	rows := u.Objects()
	groups := make([]string, len(rows))
	for i, v := range rows {
		groups[i] = v.Group
	}
	return groups, err
}

func (u *GroupPackage) MakeExpireTime(baseTime ...time.Time) time.Time {
	var baseT time.Time
	if len(baseTime) > 0 {
		baseT = baseTime[0]
	}
	if baseT.IsZero() {
		baseT = time.Now()
	}
	var t time.Time
	switch u.TimeUnit {
	case GroupPackageTimeDay:
		t = baseT.AddDate(0, 0, int(u.TimeDuration))
	case GroupPackageTimeWeek:
		t = baseT.AddDate(0, 0, int(u.TimeDuration*7))
	case GroupPackageTimeMonth:
		t = baseT.AddDate(0, int(u.TimeDuration), 0)
	case GroupPackageTimeYear:
		t = baseT.AddDate(int(u.TimeDuration), 0, 0)
	case GroupPackageTimeForever:
	}
	return t
}
