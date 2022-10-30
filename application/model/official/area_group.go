package official

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/dbschema"
)

func NewAreaGroup(ctx echo.Context) *AreaGroup {
	return &AreaGroup{
		OfficialCommonAreaGroup: dbschema.NewOfficialCommonAreaGroup(ctx),
	}
}

type AreaGroup struct {
	*dbschema.OfficialCommonAreaGroup
}

func (f *AreaGroup) Exists(countryAbbr string, abbr string) (bool, error) {
	return f.OfficialCommonAreaGroup.Exists(nil, db.And(
		db.Cond{`country_abbr`: countryAbbr},
		db.Cond{`abbr`: abbr},
	))
}

func (f *AreaGroup) ExistsOther(countryAbbr string, abbr string, id uint) (bool, error) {
	return f.OfficialCommonAreaGroup.Exists(nil, db.And(
		db.Cond{`country_abbr`: countryAbbr},
		db.Cond{`abbr`: abbr},
		db.Cond{`id`: db.NotEq(id)},
	))
}

func (f *AreaGroup) check() error {
	var (
		exists bool
		err    error
	)
	if f.Id < 1 {
		exists, err = f.Exists(f.CountryAbbr, f.Abbr)
	} else {
		exists, err = f.ExistsOther(f.CountryAbbr, f.Abbr, f.Id)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().E(`缩写“%s-%s”已存在`, f.CountryAbbr, f.Abbr)
	}

	f.CountryAbbr = strings.ToUpper(f.CountryAbbr)
	return nil
}

func (f *AreaGroup) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCommonAreaGroup.Insert()
}

func (f *AreaGroup) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCommonAreaGroup.Update(mw, args...)
}

func (f *AreaGroup) GetWithExt(cond *db.Compounds) (*AreaGroupExt, error) {
	recv := &AreaGroupExt{}
	err := f.NewParam().SetArgs(cond.And()).SetRecv(recv).One()
	return recv, err
}

func (f *AreaGroup) ListPageWithExt(cond *db.Compounds, sorts ...interface{}) ([]*AreaGroupExt, error) {
	list := []*AreaGroupExt{}
	_, err := handler.NewLister(f.OfficialCommonAreaGroup, &list, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(f.Context())
	return list, err
}

func (f *AreaGroup) ListPage(cond *db.Compounds, sorts ...interface{}) ([]*dbschema.OfficialCommonAreaGroup, error) {
	_, err := handler.NewLister(f.OfficialCommonAreaGroup, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return nil, err
	}
	return f.Objects(), err
}
