package official

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

var (
	// NavigateMaxLevel 分类最大层数(层从0开始)
	NavigateMaxLevel uint = 3 //3代表最大允许4层
)

func NewNavigate(ctx echo.Context) *Navigate {
	return &Navigate{
		OfficialCommonNavigate: dbschema.NewOfficialCommonNavigate(ctx),
		maxLevel:               NavigateMaxLevel,
	}
}

var onNavigateUpdate []func(*dbschema.OfficialCommonNavigate) error

func FireNavigateUpdate(m *dbschema.OfficialCommonNavigate) error {
	for _, f := range onNavigateUpdate {
		if f == nil {
			continue
		}
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

func OnNavigateUpdate(f func(*dbschema.OfficialCommonNavigate) error) {
	onNavigateUpdate = append(onNavigateUpdate, f)
}

type Navigate struct {
	*dbschema.OfficialCommonNavigate
	maxLevel uint
}

func (f *Navigate) MaxLevel() uint {
	return f.maxLevel
}

func (f *Navigate) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.Get(mw, args...)
	if err != nil {
		return err
	}
	if f.HasChild == `Y` {
		return f.Context().E(`删除失败：请先删除子分类`)
	}
	if err = f.Context().Begin(); err != nil {
		return err
	}
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.OfficialCommonNavigate.Delete(mw, args...)
	if err != nil {
		return err
	}
	FireNavigateUpdate(f.OfficialCommonNavigate)
	err = f.UpdateAllParents(f.OfficialCommonNavigate)
	return err
}

func (f *Navigate) ListAllParent(typ string, excludeId uint, maxLevels ...uint) []*dbschema.OfficialCommonNavigate {
	maxLevel := f.MaxLevel()
	if len(maxLevels) > 0 {
		maxLevel = maxLevels[0]
	}
	queryMW := func(r db.Result) db.Result {
		return r.OrderBy(`level`, `sort`, `id`)
	}
	cond := db.Compounds{}
	cond.AddKV(`type`, typ)
	cond.AddKV(`disabled`, `N`)
	cond.AddKV(`level`, db.Lte(maxLevel))
	if excludeId > 0 {
		cond.AddKV(`id`, db.NotEq(excludeId))
	}
	f.ListByOffset(nil, queryMW, 0, -1, cond.And())
	return f.Objects()
}

func (f *Navigate) ListIndent(categoryList []*dbschema.OfficialCommonNavigate) []*dbschema.OfficialCommonNavigate {
	for idx, row := range categoryList {
		categoryList[idx].Title = strings.Repeat(ChineseSpace, int(row.Level)) + row.Title
	}
	return categoryList
}

func (f *Navigate) ParentIds(parentID uint) []uint {
	r := []uint{}
	for parentID > 0 {
		err := f.Get(nil, `id`, parentID)
		if err != nil {
			return r
		}
		r = append(r, parentID)
		parentID = f.ParentId
	}
	result := make([]uint, len(r))
	for k, i := 0, len(r)-1; i >= 0; i-- {
		result[k] = r[i]
		k++
	}
	return result
}

func (f *Navigate) Add() (pk interface{}, err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
	if len(f.Type) == 0 {
		f.Type = `default`
	}
	if len(f.LinkType) == 0 {
		f.LinkType = `custom`
	}
	err = f.Exists(f.Title, f.Type)
	if err != nil {
		return
	}
	if f.ParentId > 0 {
		parent := dbschema.NewOfficialCommonNavigate(f.Context())
		err = parent.Get(nil, `id`, f.ParentId)
		if err != nil {
			if err != db.ErrNoMoreRows {
				return
			}
			err = f.Context().E(`父级分类不存在`)
			return
		}
		f.Level = parent.Level + 1
		f.Type = parent.Type
		if parent.HasChild == `N` {
			err = parent.UpdateField(nil, `has_child`, `Y`, `id`, f.ParentId)
			if err != nil {
				return
			}
		}
	} else {
		f.Level = 0
	}
	if f.Level > f.MaxLevel() {
		err = f.Context().E(`操作失败！菜单超过最大层数: %d`, f.MaxLevel())
		return
	}
	pk, err = f.OfficialCommonNavigate.Insert()
	if err != nil {
		return
	}
	err = FireNavigateUpdate(f.OfficialCommonNavigate)
	return
}

func (f *Navigate) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	f.Context().Begin()
	defer func() {
		if err == nil {
			err = FireNavigateUpdate(f.OfficialCommonNavigate)
		}
		f.Context().End(err == nil)
	}()
	if len(f.Type) == 0 {
		f.Type = `default`
	}
	if len(f.LinkType) == 0 {
		f.LinkType = `custom`
	}
	if err = f.ExistsOther(f.Title, f.Type, f.Id); err != nil {
		return err
	}
	oldData := dbschema.NewOfficialCommonNavigate(f.Context())
	err = oldData.Get(nil, args...)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return f.Context().E(`菜单不存在`)
	}
	if oldData.Id == f.ParentId {
		return f.Context().E(`不能选择自己为上级菜单`)
	}
	if oldData.ParentId != f.ParentId {
		if f.ParentId > 0 {
			parent := dbschema.NewOfficialCommonNavigate(f.Context())
			err = parent.Get(nil, `id`, f.ParentId)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				return f.Context().E(`父级分类不存在`)
			}
			f.Level = parent.Level + 1
			f.Type = parent.Type
			if parent.HasChild == `N` {
				err = parent.UpdateField(nil, `has_child`, `Y`, `id`, f.ParentId)
				if err != nil {
					return
				}
			}
		} else {
			f.Level = 0
		}
		if f.Level > f.MaxLevel() {
			return f.Context().E(`操作失败！分类超过最大层数: %d`, f.MaxLevel())
		}
		err = f.UpdateAllParents(oldData)
		if err != nil {
			return err
		}
		err = f.UpdateAllChildren(f.OfficialCommonNavigate)
		if err != nil {
			return err
		}
	}
	return f.OfficialCommonNavigate.Update(mw, args...)
}

// UpdateAllParents 更新所有父级分类(目前仅仅用于更新父级has_child值)
func (f *Navigate) UpdateAllParents(oldData *dbschema.OfficialCommonNavigate) (err error) {
	if oldData.ParentId == 0 {
		return
	}
	oldParent := dbschema.NewOfficialCommonNavigate(f.Context())
	err = oldParent.Get(nil, `id`, oldData.ParentId)
	if err == nil {
		ocond := db.And(
			db.Cond{`parent_id`: oldData.ParentId},
			db.Cond{`id`: db.NotEq(oldData.Id)},
		)
		var n int64
		n, err = oldData.Count(nil, ocond)
		if err != nil {
			return
		}
		if n == 0 && oldParent.HasChild != `N` {
			err = oldData.UpdateField(nil, `has_child`, `N`, `id`, oldData.ParentId)
		} else if n > 0 && oldParent.HasChild == `N` {
			err = oldData.UpdateField(nil, `has_child`, `Y`, `id`, oldData.ParentId)
		}
	}
	if err != nil && err == db.ErrNoMoreRows {
		err = nil
	}
	return
}

// UpdateAllChildren 更新所有子孙分类level值
func (f *Navigate) UpdateAllChildren(row *dbschema.OfficialCommonNavigate) error {
	children := dbschema.NewOfficialCommonNavigate(f.Context())
	_, err := children.ListByOffset(nil, nil, 0, -1, `parent_id`, row.Id)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return nil
	}
	for _, child := range children.Objects() {
		child.Level = row.Level + 1
		err = child.UpdateField(nil, `level`, child.Level, `id`, child.Id)
		if err != nil {
			return err
		}
		err = f.UpdateAllChildren(child)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *Navigate) Exists(title string, typ string) error {
	exists, err := f.OfficialCommonNavigate.Exists(nil, db.And(
		db.Cond{`title`: title},
		db.Cond{`type`: typ},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`名称“%s”已经使用过了`, title)
	}
	return err
}

func (f *Navigate) ExistsOther(title string, typ string, id uint) error {
	exists, err := f.OfficialCommonNavigate.Exists(nil, db.And(
		db.Cond{`title`: title},
		db.Cond{`type`: typ},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`名称“%s”已经使用过了`, title)
	}
	return err
}

func CollectionNavigateIDs() func(...uint) []uint {
	var categoryIds []uint
	return func(cIds ...uint) []uint {
		for _, cID := range cIds {
			if !com.InUintSlice(cID, categoryIds) {
				categoryIds = append(categoryIds, cID)
			}
		}
		return categoryIds
	}
}
