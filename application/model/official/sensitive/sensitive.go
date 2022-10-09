package sensitive

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func NewSensitive(ctx echo.Context) *Sensitive {
	return &Sensitive{
		OfficialCommonSensitive: dbschema.NewOfficialCommonSensitive(ctx),
	}
}

type Sensitive struct {
	*dbschema.OfficialCommonSensitive
}

func (f *Sensitive) Exists(words string) (bool, error) {
	return f.OfficialCommonSensitive.Exists(nil, db.Cond{`words`: words})
}

func (f *Sensitive) ExistsOther(words string, id uint) (bool, error) {
	return f.OfficialCommonSensitive.Exists(nil, db.And(
		db.Cond{`words`: words},
		db.Cond{`id <>`: id},
	))
}

func (f *Sensitive) check() error {
	return nil
}

func (f *Sensitive) Add() (pk interface{}, err error) {
	err = f.check()
	if err != nil {
		return
	}
	return f.OfficialCommonSensitive.Insert()
}

func (f *Sensitive) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.check()
	if err != nil {
		return err
	}
	return f.OfficialCommonSensitive.Update(mw, args...)
}
