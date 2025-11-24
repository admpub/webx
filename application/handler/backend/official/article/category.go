package article

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
)

func CategoryIndex(ctx echo.Context) error {
	m := official.NewCategory(ctx)
	cond := db.Compounds{}
	t := ctx.Form(`type`)
	parentID := ctx.Formx(`parentId`).Uint()
	currentID := ctx.Formx(`currentId`).Uint()
	onlyList := ctx.Formx(`onlyList`).Bool()
	if len(t) > 0 {
		cond.AddKV(`type`, t)
	}
	cond.AddKV(`parent_id`, parentID)
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

// CategoryAdd TODO: 多语种支持
func CategoryAdd(ctx echo.Context) error {
	var err error
	m := official.NewCategory(ctx)
	t := ctx.Form(`type`, `article`)
	parentID := ctx.Formx(`parentId`).Uint()
	if parentID > 0 {
		err = m.Get(nil, db.Cond{`id`: parentID})
		if err != nil {
			return err
		}
		t = m.Type
		m.HasChild = `Y`
	}
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(m.OfficialCommonCategory, uint64(id))
			} else {
				m.Sort = 5000
			}
		} else {
			if parentID > 0 {
				m.ParentId = parentID
			}
			m.Sort = 5000
		}
	}
	if len(m.Type) == 0 {
		m.Type = t
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonCategory,
		formbuilder.ConfigFile(`official/article/category_edit`),
		formbuilder.AllowedNames(
			`type`, `parentId`, `cover`, `name`, `keywords`, `slugify`,
			`description`, `template`, `showOnMenu`, `sort`, `disabled`,
		),
	)
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(m.OfficialCommonCategory, uint64(m.Id))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(backend.URLFor(`/official/article/category`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()

	ctx.Set(`activeURL`, `/official/article/category`)
	categoryList := m.ListIndent(m.ListAllParent(m.Type, 0))
	ctx.Set(`categoryList`, categoryList)
	ctx.Set(`typeList`, official.CategoryTypes.Slice())
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

// CategoryEdit TODO: 多语种支持
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
	allowedNames := []string{`parentId`, `cover`, `name`, `keywords`, `slugify`,
		`description`, `template`, `showOnMenu`, `sort`, `disabled`}
	if editableType {
		allowedNames = append(allowedNames, `type`)
	}
	if ctx.IsGet() {
		i18nm.SetModelTranslationsToForm(m.OfficialCommonCategory, uint64(id))
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonCategory,
		formbuilder.ConfigFile(`official/article/category_edit`),
		formbuilder.AllowedNames(allowedNames...),
	)
	form.OnPost(func() error {
		err = m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(m.OfficialCommonCategory, uint64(m.Id))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/article/category`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
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
