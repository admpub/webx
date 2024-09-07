package advert

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
)

func NewAdPosition(ctx echo.Context) *AdPosition {
	a := &AdPosition{
		OfficialAdPosition: dbschema.NewOfficialAdPosition(ctx),
	}
	return a
}

type AdPosition struct {
	*dbschema.OfficialAdPosition
}

func (f *AdPosition) Exists(name string, ident string, excludeId uint64) error {
	bean := dbschema.NewOfficialAdPosition(f.Context())
	cond := db.NewCompounds()
	cond.Add(db.Or(
		db.Cond{`name`: name},
		db.Cond{`ident`: ident},
	))
	if excludeId > 0 {
		cond.AddKV(`id`, db.NotEq(excludeId))
	}
	err := bean.Get(nil, cond.And())
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	if bean.Name == name {
		return f.Context().NewError(code.DataAlreadyExists, `广告位名称“%s”已经存在`, name).SetZone(`name`)
	}
	return f.Context().NewError(code.DataAlreadyExists, `广告位唯一标识“%s”已经存在`, ident).SetZone(`ident`)
}

func (f *AdPosition) check() error {
	if len(f.Name) == 0 {
		return f.Context().NewError(code.InvalidParameter, `广告位名称不能为空`, f.Name).SetZone(`name`)
	}
	if len(f.Ident) == 0 {
		return f.Context().NewError(code.InvalidParameter, `广告位标识不能为空`, f.Ident).SetZone(`ident`)
	}
	if strings.Contains(f.Ident, `,`) {
		return f.Context().NewError(code.InvalidParameter, `禁止在广告位标识“%s”中使用“,”，请修改`, f.Ident).SetZone(`ident`)
	}
	err := f.Exists(f.Name, f.Ident, f.Id)
	if err != nil {
		return err
	}

	return nil
}

func (f *AdPosition) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialAdPosition.Insert()
}

func (f *AdPosition) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialAdPosition.Update(mw, args...)
}

func (f *AdPosition) GetIDByIdent(ident string) (uint64, error) {
	err := f.OfficialAdPosition.Get(func(r db.Result) db.Result {
		return r.Select(`id`)
	}, db.And(
		db.Cond{`ident`: ident},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
		return 0, err
	}
	return f.Id, err
}

var ErrAdSlotDoesNotExist = errors.New(`ad slot does not exist`)

func (f *AdPosition) GetOneAdvertByIdent(ident string) (*dbschema.OfficialAdItem, error) {
	posID, err := f.GetIDByIdent(ident)
	if err != nil {
		return nil, err
	}
	if posID == 0 {
		return nil, fmt.Errorf(`%w: %s`, ErrAdSlotDoesNotExist, ident)
	}
	m := dbschema.NewOfficialAdItem(f.Context())
	cond := db.NewCompounds()
	cond.Add(db.Cond{`position_id`: posID})
	err = m.Get(func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, cond.And(GenAdvertCondition()...))
	return m, err
}

func GenAdvertCondition() db.Compounds {
	nowTS := time.Now().Unix()
	return []db.Compound{
		db.Cond{`disabled`: `N`},
		db.Cond{`start`: db.Lte(nowTS)},
		db.Or(
			db.Cond{`end`: db.Gte(nowTS)},
			db.Cond{`end`: 0},
		),
	}
}

func (f *AdPosition) GetAdvertsByIdent(idents ...string) (PositionAdverts, error) {
	_, err := f.OfficialAdPosition.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`ident`: db.In(idents)},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
		return nil, err
	}
	positions := f.OfficialAdPosition.Objects()
	if len(positions) == 0 {
		return nil, nil
	}
	posIDs := make([]uint64, len(positions))
	posIDi := map[uint64]*dbschema.OfficialAdPosition{}
	list := PositionAdverts{}
	for idx, pos := range positions {
		posIDs[idx] = pos.Id
		posIDi[pos.Id] = pos
	}
	item := dbschema.NewOfficialAdItem(f.Context())
	cond := db.NewCompounds()
	cond.Add(db.Cond{`position_id`: db.In(posIDs)})
	_, err = item.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`position_id`, `sort`, `id`)
	}, 0, -1, cond.And(GenAdvertCondition()...))
	for _, row := range item.Objects() {
		pos := posIDi[row.PositionId]
		if _, ok := list[pos.Ident]; !ok {
			list[pos.Ident] = ItemsResponse{}
		}
		list[pos.Ident] = append(list[pos.Ident], NewItemResponse(row, pos))
	}
	for _, pos := range positions {
		if _, ok := list[pos.Ident]; ok {
			continue
		}
		item := &dbschema.OfficialAdItem{
			Content: pos.Content,
			Contype: pos.Contype,
			Url:     pos.Url,
		}
		list[pos.Ident] = []*ItemResponse{NewItemResponse(item, pos)}
	}
	return list, err
}
