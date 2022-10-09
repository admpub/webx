package advert

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
)

func NewAdPublisher(ctx echo.Context) *AdPublisher {
	a := &AdPublisher{
		OfficialAdPublisher: dbschema.NewOfficialAdPublisher(ctx),
	}
	return a
}

type AdPublisher struct {
	*dbschema.OfficialAdPublisher
}

func (f *AdPublisher) Exists(ownerType string, ownerID uint64) (bool, error) {
	return f.OfficialAdPublisher.Exists(nil, db.And(
		db.Cond{`owner_type`: ownerType},
		db.Cond{`owner_id`: ownerID},
	))
}

func (f *AdPublisher) ExistsOther(ownerType string, ownerID uint64, id uint64) (bool, error) {
	return f.OfficialAdPublisher.Exists(nil, db.And(
		db.Cond{`owner_type`: ownerType},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`id`: db.NotEq(id)},
	))
}

func (f *AdPublisher) check() error {
	var (
		exists bool
		err    error
	)
	if f.Id < 1 {
		exists, err = f.Exists(f.OwnerType, f.OwnerId)
	} else {
		exists, err = f.ExistsOther(f.OwnerType, f.OwnerId, f.Id)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().NewError(code.DataAlreadyExists, `用户“%s: %d”已存在`, f.OwnerType, f.OwnerId)
	}

	return nil
}

func (f *AdPublisher) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialAdPublisher.Insert()
}

func (f *AdPublisher) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialAdPublisher.Update(mw, args...)
}
