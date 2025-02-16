package article

import (
	"strconv"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/official"
)

func CategoryIndex(ctx echo.Context) error {
	m := official.NewCategory(ctx)
	cond := db.Compounds{}
	t := ctx.Form(`type`)
	parentID := ctx.Formx(`parentId`).Int()
	currentID := ctx.Formx(`currentId`).Uint()
	onlyList := ctx.Formx(`onlyList`).Bool()
	if len(t) > 0 {
		cond.AddKV(`type`, t)
	}
	if parentID >= 0 {
		cond.AddKV(`parent_id`, parentID)
	}
	if currentID > 0 {
		cond.AddKV(`id`, db.NotEq(currentID))
	}
	queryMW := func(r db.Result) db.Result {
		return r.OrderBy(`type`, `level`, `sort`, `id`)
	}
	if onlyList {
		_, err := m.ListByOffset(nil, queryMW, 0, -1, cond.And())
		if err != nil {
			return err
		}
		data := ctx.Data()
		data.SetData(m.Objects())
		return ctx.JSON(data)
	}
	_, err := common.NewLister(m.OfficialCommonCategory, nil, queryMW, cond.And()).Paging(ctx)
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`maxLevel`, official.CategoryMaxLevel)
	ctx.SetFunc(`getTypeName`, official.CategoryTypes.Get)
	if ctx.Formx(`partial`).Bool() {
		return ctx.Render(`official/article/category_list_row`, ret)
	}
	ctx.Set(`typeList`, official.CategoryTypes.Slice())
	ctx.Set(`type`, t)
	return ctx.Render(`official/article/category`, ret)
}

func categoryEdiableType(ctx echo.Context, m *dbschema.OfficialCommonCategory) bool {
	var editable bool
	if m == nil || m.Id < 1 {
		editable = true
	} else {
		editable = m.ParentId == 0 && m.HasChild == `N`
	}
	ctx.Set(`editableType`, editable)
	return editable
}

func CategoryAdd(ctx echo.Context) error {
	var err error
	m := official.NewCategory(ctx)
	t := ctx.Form(`type`, `article`)
	parentID := ctx.Formx(`parentId`).Int()
	if parentID > 0 {
		err = m.Get(nil, db.Cond{`id`: parentID})
		if err != nil {
			return err
		}
		t = m.Type
		m.HasChild = `Y`
	}
	if ctx.IsPost() {
		m.Reset()
		err = ctx.MustBind(m.OfficialCommonCategory)
		if err == nil {
			var added []string
			var first bool
			added, err = common.BatchAdd(ctx, `name,slugify`, m, func(v *string) error {
				m.Id = 0
				m.HasChild = `N`
				if !first {
					first = true
				} else {
					m.Slugify = ``
				}
				return nil
			}, `=`)
			if err == nil && len(added) == 0 {
				err = ctx.E(`分类名称不能为空`)
			}
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/article/category`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonCategory, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		} else {
			if parentID > 0 {
				ctx.Request().Form().Set(`parentId`, strconv.Itoa(parentID))
			}
		}
	}

	ctx.Set(`activeURL`, `/official/article/category`)
	categoryList := m.ListIndent(m.ListAllParent(t, 0))
	ctx.Set(`categoryList`, categoryList)
	ctx.Set(`typeList`, official.CategoryTypes.Slice())
	ctx.Set(`type`, t)
	categoryEdiableType(ctx, m.OfficialCommonCategory)
	ctx.SetFunc(`getTypeName`, official.CategoryTypes.Get)
	return ctx.Render(`official/article/category_edit`, err)
}

func ajaxCategorySetDisabled(ctx echo.Context) error {
	return ajaxCategorySetSwitch(ctx, `disabled`, `disabled`)
}

func ajaxCategorySetShowOnMenu(ctx echo.Context) error {
	return ajaxCategorySetSwitch(ctx, `showOnMenu`, `show_on_menu`)
}

func ajaxCategorySetSwitch(ctx echo.Context, inputKey, field string) error {
	id := ctx.Formx(`id`).Uint()
	m := official.NewCategory(ctx)
	value := ctx.Query(inputKey)
	if !common.IsBoolFlag(value) {
		return ctx.NewError(code.InvalidParameter, ``).SetZone(inputKey)
	}
	data := ctx.Data()
	err := m.UpdateField(nil, field, value, db.Cond{`id`: id})
	if err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	data.SetInfo(ctx.T(`操作成功`))
	return ctx.JSON(data)
}

func ajaxCategorySetSort(ctx echo.Context) error {
	id := ctx.Formx(`pk`).Uint()
	m := official.NewCategory(ctx)
	sort := ctx.Formx(`value`).Int()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, ``).SetZone(`pk`)
	}
	data := ctx.Data()
	err := m.UpdateField(nil, `sort`, sort, db.Cond{`id`: id})
	if err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	data.SetInfo(ctx.T(`操作成功`))
	return ctx.JSON(data)
}

var ajaxCategorySet = echo.HandlerFuncs{
	`setDisabled`:   ajaxCategorySetDisabled,
	`setShowOnMenu`: ajaxCategorySetShowOnMenu,
	`setSort`:       ajaxCategorySetSort,
}

func CategoryEdit(ctx echo.Context) error {
	op := ctx.Form(`op`)
	if len(op) > 0 {
		return ajaxCategorySet.Call(ctx, op)
	}
	var err error
	id := ctx.Formx(`id`).Uint()
	m := official.NewCategory(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	editableType := categoryEdiableType(ctx, m.OfficialCommonCategory)
	if editableType {
		t := ctx.Form(`type`)
		if len(t) > 0 {
			m.Type = t
		}
	}
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`分类名称不能为空`)
		} else if e := m.ExistsOther(name, id); e != nil {
			err = e
		} else {
			excludeFields := []string{`created`}
			if editableType {
				excludeFields = append(excludeFields, `type`)
			}
			err = ctx.MustBind(m.OfficialCommonCategory, echo.ExcludeFieldName(excludeFields...))
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/article/category`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonCategory, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/official/article/category`)
	categoryRows := m.ListAllParent(m.Type, 0)
	categoryList := []*dbschema.OfficialCommonCategory{}
	var (
		exCurrentLevel int64 = -1
	)
	for _, row := range categoryRows {
		if row.Id == m.Id {
			exCurrentLevel = int64(row.Level)
			continue
		}
		if exCurrentLevel >= 0 {
			if int64(row.Level) == exCurrentLevel {
				exCurrentLevel = -1
			} else {
				continue
			}
		}
		row.Name = strings.Repeat(official.ChineseSpace, int(row.Level)) + row.Name
		categoryList = append(categoryList, row)
	}
	ctx.Set(`categoryList`, categoryList)
	ctx.Set(`typeList`, official.CategoryTypes.Slice())
	ctx.Set(`type`, m.Type)
	ctx.SetFunc(`getTypeName`, official.CategoryTypes.Get)
	return ctx.Render(`official/article/category_edit`, err)
}

func CategoryDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := official.NewCategory(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/article/category`))
}
