package advert

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func NewAdSettings(ctx echo.Context) *AdSettings {
	a := &AdSettings{
		OfficialAdSettings: dbschema.NewOfficialAdSettings(ctx),
	}
	return a
}

type AdSettings struct {
	*dbschema.OfficialAdSettings
}

func (f *AdSettings) Exists(name string) (bool, error) {
	return f.OfficialAdSettings.Exists(nil, db.Cond{`name`: name})
}

func (f *AdSettings) ExistsOther(name string, id uint64) (bool, error) {
	return f.OfficialAdSettings.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}

func (f *AdSettings) check() error {
	return nil
}

func (f *AdSettings) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialAdSettings.Insert()
}

func (f *AdSettings) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialAdSettings.Update(mw, args...)
}

// TODO: unimplemented
func (f *AdSettings) List(filter Filter, offset, size int) error {
	cond := filter.GenCond()
	_, err := f.OfficialAdSettings.ListByOffset(nil, nil, offset, size, cond.And())
	if err != nil {
		return err
	}
	return err
}
