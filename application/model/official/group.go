package official

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func NewGroup(ctx echo.Context) *Group {
	return &Group{
		OfficialCommonGroup: dbschema.NewOfficialCommonGroup(ctx),
	}
}

type Group struct {
	*dbschema.OfficialCommonGroup
}

func (f *Group) Exists(name string) (bool, error) {
	return f.OfficialCommonGroup.Exists(nil, db.Cond{`name`: name})
}

func (f *Group) ExistsOther(name string, id uint) (bool, error) {
	return f.OfficialCommonGroup.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}

func (f *Group) Add() (pk interface{}, err error) {
	return f.OfficialCommonGroup.Insert()
}

func (f *Group) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonGroup.Update(mw, args...)
}
