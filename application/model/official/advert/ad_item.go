package advert

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
)

func NewAdItem(ctx echo.Context) *AdItem {
	a := &AdItem{
		OfficialAdItem: dbschema.NewOfficialAdItem(ctx),
	}
	return a
}

type AdItem struct {
	*dbschema.OfficialAdItem
}

func (f *AdItem) Exists(name string) (bool, error) {
	return f.OfficialAdItem.Exists(nil, db.Cond{`name`: name})
}

func (f *AdItem) ExistsOther(name string, id uint64) (bool, error) {
	return f.OfficialAdItem.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}

func (f *AdItem) check() error {
	if len(f.Name) == 0 {
		return f.Context().NewError(code.InvalidParameter, `广告名称不能为空`, f.Name).SetZone(`name`)
	}
	var (
		exists bool
		err    error
	)
	if f.Id < 1 {
		exists, err = f.Exists(f.Name)
	} else {
		exists, err = f.ExistsOther(f.Name, f.Id)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().E(`名称“%s”已存在`, f.Name)
	}

	return nil
}

func (f *AdItem) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialAdItem.Insert()
}

func (f *AdItem) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialAdItem.Update(mw, args...)
}
