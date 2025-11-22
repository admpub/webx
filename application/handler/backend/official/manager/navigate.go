package manager

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/official"
)

func NavigateIndex(ctx echo.Context) error {
	m := official.NewNavigate(ctx)
	cond := db.Compounds{}
	t := ctx.Form(`type`)
	parentID := ctx.Formx(`parentId`).Uint()
	currentID := ctx.Formx(`currentId`).Uint()
	onlyList := ctx.Formx(`onlyList`).Bool()
	if len(t) > 0 {
		cond.AddKV(`type`, t)
	}
	if parentID > 0 {
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
	_, err := common.NewLister(m.OfficialCommonNavigate, nil, queryMW, cond.And()).Paging(ctx)
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`maxLevel`, official.NavigateMaxLevel)
	ctx.SetFunc(`getTypeName`, official.NavigateTypes.Get)
	ctx.SetFunc(`getLinkTypeName`, official.NavigateLinkType.Get)
	if ctx.Formx(`partial`).Bool() {
		return ctx.Render(`official/manager/navigate/list_row`, ret)
	}
	ctx.Set(`typeList`, official.GetAddNavigateTypes())
	ctx.Set(`type`, t)
	return ctx.Render(`official/manager/navigate/index`, ret)
}

func navigateEdiableType(ctx echo.Context, m *dbschema.OfficialCommonNavigate) bool {
	var editable bool
	if m == nil || m.Id < 1 {
		editable = true
	} else {
		editable = m.ParentId == 0 && m.HasChild == `N`
	}
	ctx.Set(`editableType`, editable)
	return editable
}

func NavigateAdd(ctx echo.Context) error {
	var err error
	m := official.NewNavigate(ctx)
	t := ctx.Form(`type`, `default`)
	parentID := ctx.Formx(`parentId`).Uint()
	if parentID > 0 {
		err = m.Get(nil, db.Cond{`id`: parentID})
		if err != nil {
			return err
		}
		t = m.Type
		m.HasChild = `Y`
	}

	form := formbuilder.New(ctx,
		m.OfficialCommonNavigate,
		formbuilder.ConfigFile(`official/manager/navigate/edit`))
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(backend.URLFor(`/manager/navigate/index`))
	})
	form.OnGet(func() error {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonNavigate, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		} else {
			if parentID > 0 {
				ctx.Request().Form().Set(`parentId`, param.AsString(parentID))
			}
		}
		return nil
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	/*
		if ctx.IsPost() {
			m.Reset()
			err = ctx.MustBind(m.OfficialCommonNavigate)
			if err == nil {
				var added []string
				added, err = common.BatchAdd(ctx, `title,ident,url`, m, func(_ *string) error {
					m.Id = 0
					m.HasChild = `N`
					return nil
				}, `|`)
				if err == nil && len(added) == 0 {
					err = ctx.E(`菜单标题不能为空`)
				}
			}
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/manager/navigate/index`))
			}
		} else {
			id := ctx.Formx(`copyId`).Uint()
			if id > 0 {
				err = m.Get(nil, `id`, id)
				if err == nil {
					echo.StructToForm(ctx, m.OfficialCommonNavigate, ``, echo.LowerCaseFirstLetter)
					ctx.Request().Form().Set(`id`, `0`)
				}
			} else {
				if parentID > 0 {
					ctx.Request().Form().Set(`parentId`, param.AsString(parentID))
				}
			}
		}
	*/
	ctx.Set(`activeURL`, `/manager/navigate/index`)
	navigateList := m.ListIndent(m.ListAllParent(t, 0))
	ctx.Set(`navigateList`, navigateList)
	ctx.Set(`typeList`, official.GetAddNavigateTypes())
	ctx.Set(`linkTypeList`, official.NavigateLinkType.Slice())
	ctx.Set(`type`, t)
	ctx.Set(`title`, ctx.T(`添加菜单`))
	navigateEdiableType(ctx, m.OfficialCommonNavigate)
	ctx.SetFunc(`getTypeName`, official.NavigateTypes.Get)
	return ctx.Render(`official/manager/navigate/edit`, err)
}

func ajaxNavigateSetDisabled(ctx echo.Context) error {
	return ajaxNavigateSetSwitch(ctx, `disabled`, `disabled`)
}

func ajaxNavigateSetSwitch(ctx echo.Context, inputKey, field string) error {
	id := ctx.Formx(`id`).Uint()
	m := official.NewNavigate(ctx)
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

func ajaxNavigateSetSort(ctx echo.Context) error {
	id := ctx.Formx(`pk`).Uint()
	m := official.NewNavigate(ctx)
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

var ajaxNavigateSet = echo.HandlerFuncs{
	`setDisabled`: ajaxNavigateSetDisabled,
	`setSort`:     ajaxNavigateSetSort,
}

func NavigateEdit(ctx echo.Context) error {
	op := ctx.Form(`op`)
	if len(op) > 0 {
		return ajaxNavigateSet.Call(ctx, op)
	}
	var err error
	id := ctx.Formx(`id`).Uint()
	m := official.NewNavigate(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	editableType := navigateEdiableType(ctx, m.OfficialCommonNavigate)
	if editableType {
		t := ctx.Form(`type`)
		if len(t) > 0 {
			m.Type = t
		}
	}
	if ctx.IsPost() {
		excludeFields := []string{`created`}
		if editableType {
			excludeFields = append(excludeFields, `type`)
		}
		err = ctx.MustBind(m.OfficialCommonNavigate, echo.ExcludeFieldName(excludeFields...))

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/manager/navigate/index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonNavigate, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/manager/navigate`)
	navigateRows := m.ListAllParent(m.Type, 0)
	navigateList := []*dbschema.OfficialCommonNavigate{}
	var (
		exCurrentLevel int64 = -1
	)
	for _, row := range navigateRows {
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
		row.Title = strings.Repeat(official.ChineseSpace, int(row.Level)) + row.Title
		navigateList = append(navigateList, row)
	}
	ctx.Set(`navigateList`, navigateList)
	ctx.Set(`typeList`, official.GetAddNavigateTypes())
	ctx.Set(`linkTypeList`, official.NavigateLinkType.Slice())
	ctx.Set(`type`, m.Type)
	ctx.Set(`title`, ctx.T(`修改菜单`))
	ctx.SetFunc(`getTypeName`, official.NavigateTypes.Get)
	return ctx.Render(`official/manager/navigate/edit`, err)
}

func NavigateDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := official.NewNavigate(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/manager/navigate/index`))
}
