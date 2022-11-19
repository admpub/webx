package official

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/top"
)

var (
	// CategoryMaxLevel 分类最大层数(层从0开始)
	CategoryMaxLevel uint = 3 //3代表最大允许4层
	// CategoryTypes 分类的类别
	CategoryTypes = echo.NewKVData()
)

func init() {
	CategoryTypes.AddItem(&echo.KV{K: `article`, V: `文章`})
	CategoryTypes.AddItem(&echo.KV{K: `friendlink`, V: `友情链接`})
}

// AddCategoryType 登记新的类别
func AddCategoryType(value, label string) {
	CategoryTypes.Add(value, label)
}

// ChineseSpace 中文全角空白字符
const ChineseSpace = `　`

func NewCategory(ctx echo.Context) *Category {
	m := &Category{
		OfficialCommonCategory: dbschema.NewOfficialCommonCategory(ctx),
		maxLevel:               CategoryMaxLevel,
	}
	return m
}

type Category struct {
	*dbschema.OfficialCommonCategory
	maxLevel uint
}

func (f *Category) MaxLevel() uint {
	return f.maxLevel
}

func (f *Category) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
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
	err = f.OfficialCommonCategory.Delete(mw, args...)
	if err != nil {
		return err
	}
	err = f.UpdateAllParents(f.OfficialCommonCategory)
	return err
}

func (f *Category) ListAllParent(typ string, excludeId uint, maxLevels ...uint) []*dbschema.OfficialCommonCategory {
	maxLevel := f.MaxLevel()
	if len(maxLevels) > 0 {
		maxLevel = maxLevels[0]
	}
	return f.ListAllParentBy(typ, excludeId, maxLevel)
}

func (f *Category) ListAllParentBy(typ string, excludeId uint, maxLevel uint, extraConds ...db.Compound) []*dbschema.OfficialCommonCategory {
	queryMW := func(r db.Result) db.Result {
		return r.OrderBy(`level`, `parent_id`, `sort`, `id`)
	}
	cond := db.NewCompounds()
	cond.AddKV(`type`, typ)
	cond.AddKV(`disabled`, `N`)
	cond.AddKV(`level`, db.Lte(maxLevel))
	if excludeId > 0 {
		cond.AddKV(`id`, db.NotEq(excludeId))
	}
	if len(extraConds) > 0 {
		cond.Add(extraConds...)
	}
	f.ListByOffset(nil, queryMW, 0, -1, cond.And())
	return f.Objects()
}

func (f *Category) ListIndent(categoryList []*dbschema.OfficialCommonCategory) []*dbschema.OfficialCommonCategory {
	for idx, row := range categoryList {
		categoryList[idx].Name = strings.Repeat(ChineseSpace, int(row.Level)) + row.Name
	}
	return categoryList
}

func (f *Category) Parents(parentID uint) ([]dbschema.OfficialCommonCategory, error) {
	categories := map[uint]dbschema.OfficialCommonCategory{}
	var r []uint
	for parentID > 0 && !com.InUintSlice(parentID, r) {
		err := f.Get(nil, `id`, parentID)
		if err != nil {
			if err == db.ErrNoMoreRows {
				break
			}
			return nil, err
		}
		categories[parentID] = *f.OfficialCommonCategory
		r = append(r, parentID)
		parentID = f.ParentId
	}
	result := make([]dbschema.OfficialCommonCategory, len(r))
	for k, i := 0, len(r)-1; i >= 0; i-- {
		result[k] = categories[r[i]]
		k++
	}
	return result, nil
}

func (f *Category) ParentIds(parentID uint) []uint {
	var r []uint
	for parentID > 0 && !com.InUintSlice(parentID, r) {
		err := f.Get(nil, `id`, parentID)
		if err != nil {
			if err == db.ErrNoMoreRows {
				break
			}
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

func (f *Category) Add() (pk interface{}, err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.Exists(f.Name)
	if err != nil {
		return
	}
	if f.ParentId > 0 {
		parent := dbschema.NewOfficialCommonCategory(f.Context())
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
		err = f.Context().E(`操作失败！分类超过最大层数: %d`, f.MaxLevel())
		return
	}
	f.Slugify = top.Slugify(f.Name)
	return f.OfficialCommonCategory.Insert()
}

func (f *Category) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
	if err = f.ExistsOther(f.Name, f.Id); err != nil {
		return err
	}
	oldData := dbschema.NewOfficialCommonCategory(f.Context())
	err = oldData.Get(nil, args...)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return f.Context().E(`分类不存在`)
	}
	if oldData.Id == f.ParentId {
		return f.Context().E(`不能选择自己为上级分类`)
	}
	if oldData.ParentId != f.ParentId {
		if f.ParentId > 0 {
			parent := dbschema.NewOfficialCommonCategory(f.Context())
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
		err = f.UpdateAllChildren(f.OfficialCommonCategory)
		if err != nil {
			return err
		}
	}
	f.Slugify = top.Slugify(f.Name)
	return f.OfficialCommonCategory.Update(mw, args...)
}

// UpdateAllParents 更新所有父级分类(目前仅仅用于更新父级has_child值)
func (f *Category) UpdateAllParents(oldData *dbschema.OfficialCommonCategory) (err error) {
	if oldData.ParentId == 0 {
		return
	}
	oldParent := dbschema.NewOfficialCommonCategory(f.Context())
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
func (f *Category) UpdateAllChildren(row *dbschema.OfficialCommonCategory) error {
	children := dbschema.NewOfficialCommonCategory(f.Context())
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

func (f *Category) Exists(name string) error {
	exists, err := f.OfficialCommonCategory.Exists(nil, db.Cond{`name`: name})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`分类名称“%s”已经使用过了`, name)
	}
	return err
}

func (f *Category) ExistsOther(name string, id uint) error {
	exists, err := f.OfficialCommonCategory.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`分类名称“%s”已经使用过了`, name)
	}
	return err
}

func (f *Category) FillTo(tg []ICategory) error {
	var categoryIds []uint
	for _, u := range tg {
		if u.GetCategory1() > 0 {
			if !com.InUintSlice(u.GetCategory1(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategory1())
			}
		}
		if u.GetCategory2() > 0 {
			if !com.InUintSlice(u.GetCategory2(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategory2())
			}
		}
		if u.GetCategory3() > 0 {
			if !com.InUintSlice(u.GetCategory3(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategory3())
			}
		}
		if u.GetCategoryID() > 0 {
			if !com.InUintSlice(u.GetCategoryID(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategoryID())
			}
		}
	}
	_, err := f.ListByOffset(nil, nil, 0, -1, db.Cond{`id IN`: categoryIds})
	if err != nil {
		return err
	}
	categoryList := f.Objects()
	categoryMap := map[uint]*dbschema.OfficialCommonCategory{}
	for _, g := range categoryList {
		categoryMap[g.Id] = g
	}
	for _, v := range tg {
		if v.GetCategory1() > 0 {
			if g, y := categoryMap[v.GetCategory1()]; y {
				v.AddCategory(g)
			}
		}
		if v.GetCategory2() > 0 {
			if g, y := categoryMap[v.GetCategory2()]; y {
				v.AddCategory(g)
			}
		}
		if v.GetCategory3() > 0 {
			if g, y := categoryMap[v.GetCategory3()]; y {
				v.AddCategory(g)
			}
		}
		if v.GetCategoryID() > 0 && (v.GetCategoryID() != v.GetCategory1() && v.GetCategoryID() != v.GetCategory2() && v.GetCategoryID() != v.GetCategory3()) {
			if g, y := categoryMap[v.GetCategoryID()]; y {
				v.AddCategory(g)
			}
		}
	}
	return nil
}

func CollectionCategoryIDs() func(...uint) []uint {
	categoryIds := []uint{}
	return func(cIds ...uint) []uint {
		for _, cID := range cIds {
			if !com.InUintSlice(cID, categoryIds) {
				categoryIds = append(categoryIds, cID)
			}
		}
		return categoryIds
	}
}
