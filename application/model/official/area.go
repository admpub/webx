package official

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/top"
)

func NewArea(ctx echo.Context) *Area {
	return &Area{
		OfficialCommonArea: dbschema.NewOfficialCommonArea(ctx),
	}
}

type Area struct {
	*dbschema.OfficialCommonArea
}

func (f *Area) Exists(name string) (bool, error) {
	return f.OfficialCommonArea.Exists(nil, db.Cond{`name`: name})
}

func (f *Area) ExistsOther(name string, id uint) (bool, error) {
	return f.OfficialCommonArea.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}

func (f *Area) check() error {
	ctx := f.Context()
	f.Name = strings.TrimSpace(f.Name)
	if len(f.Name) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入地区名称`).SetZone(`name`)
	}
	f.Short = strings.TrimSpace(f.Short)
	if len(f.Short) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入地区简称`).SetZone(`short`)
	}
	if len(f.CountryAbbr) != 2 || !com.StrIsAlpha(f.CountryAbbr) {
		return ctx.NewError(code.InvalidParameter, `请输入两个字母的国家码`).SetZone(`countryAbbr`)
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
		return ctx.NewError(code.DataAlreadyExists, `名称“%s”已存在`, f.Name).SetZone(`name`)
	}
	f.Merged = f.Name
	f.Level = 1
	if f.Pid > 0 {
		if f.Pid == f.Id {
			return f.Context().E(`不能选择当前地区数据作为上级地区`)
		}
		positions, err := f.Positions(f.Pid)
		if err != nil {
			return err
		}
		areas := make([]string, len(positions))
		for i, v := range positions {
			areas[i] = v.Name
		}
		areas = append(areas, f.Name)
		f.Merged = strings.Join(areas, `,`)
		f.Level = uint(len(areas))
	}
	f.Pinyin = strings.TrimSpace(f.Pinyin)
	if len(f.Pinyin) == 0 {
		f.Pinyin = top.Pinyin(f.Short)
	} else {
		if !com.StrIsAlpha(f.Pinyin) {
			return ctx.NewError(code.InvalidParameter, `请输入拼音字母不正确`).SetZone(`pinyin`)
		}
		f.Pinyin = strings.ToLower(f.Pinyin)
	}
	if len(f.Pinyin) > 1 {
		f.First = strings.ToUpper(f.Pinyin[0:1])
	} else if len(f.Pinyin) == 1 {
		f.First = strings.ToUpper(f.Pinyin)
	} else {
		f.First = ``
	}
	f.CountryAbbr = strings.ToUpper(f.CountryAbbr)
	return nil
}

func (f *Area) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCommonArea.Insert()
}

func (f *Area) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCommonArea.Update(mw, args...)
}

func (f *Area) Parent(pid uint) (*dbschema.OfficialCommonArea, error) {
	m := dbschema.NewOfficialCommonArea(f.Context())
	err := m.Get(nil, `id`, pid)
	return m, err
}

func (f *Area) Parents(pid uint) ([]*dbschema.OfficialCommonArea, error) {
	parents := make([]*dbschema.OfficialCommonArea, 0)
	var (
		m    *dbschema.OfficialCommonArea
		err  error
		pids = map[uint]struct{}{}
	)

loop:
	if pid == 0 {
		return parents, nil
	}
	if _, ok := pids[pid]; ok {
		return parents, nil
	}

	m, err = f.Parent(pid)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return parents, nil
		}
		return parents, err
	}
	pids[pid] = struct{}{}
	pid = m.Pid
	parents = append(parents, m)
	goto loop
}

func (f *Area) Positions(id uint) ([]*dbschema.OfficialCommonArea, error) {
	parents, err := f.Parents(id)
	if err != nil {
		return parents, err
	}
	if len(parents) == 0 {
		return parents, nil
	}
	positions := make([]*dbschema.OfficialCommonArea, len(parents))
	var index int
	for end := len(parents) - 1; end >= 0; end-- {
		positions[index] = parents[end]
		index++
	}
	return positions, err
}

func (f *Area) PositionIds(id uint) ([]uint, error) {
	parents, err := f.Parents(id)
	if err != nil {
		return nil, err
	}
	if len(parents) == 0 {
		return nil, nil
	}
	ids := make([]uint, len(parents))
	var index int
	for end := len(parents) - 1; end >= 0; end-- {
		ids[index] = parents[end].Id
		index++
	}
	return ids, err
}
