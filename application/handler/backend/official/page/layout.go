package page

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webfront/dbschema"
	modelPage "github.com/coscms/webfront/model/official/page"
)

func LayoutIndex(ctx echo.Context) error {
	m := modelPage.NewLayout(ctx)
	cond := db.NewCompounds()
	pageID := ctx.Formx(`pageId`).Uint()
	nsql.SelectPageCond(ctx, cond, `id`, `name%`)
	var pageData *dbschema.OfficialPage
	if pageID > 0 {
		pageM := modelPage.New(ctx)
		err := pageM.Get(nil, db.Cond{`id`: pageID})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return ctx.E(`页面不存在`)
			}
			return err
		}
		pageData = pageM.OfficialPage
		cond.AddKV(`page_id`, pageID)
	}
	list, err := m.ListPage(cond, `-id`)
	ctx.Set(`listData`, list)
	ctx.Set(`pageData`, pageData)
	ctx.Set(`activeURL`, `/official/page/index`)
	return ctx.Render(`official/page/layout_index`, common.Err(ctx, err))
}

func LayoutAdd(ctx echo.Context) error {
	var err error
	m := modelPage.NewLayout(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialPageLayout, echo.ExcludeFieldName(`updated`, `created`))
		if err == nil {
			_, err = m.Insert()
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/page/layout_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialPageLayout, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/page/layout_index`)
	ctx.Set(`title`, ctx.T(`添加布局`))
	return ctx.Render(`official/page/layout_edit`, common.Err(ctx, err))
}

func LayoutEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := modelPage.NewLayout(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialPageLayout, echo.ExcludeFieldName(`updated`, `created`))
		if err == nil {
			m.Id = id
			err = m.Update(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/page/layout_index`))
			}
		}
	} else if ctx.IsAjax() {

		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			data := ctx.Data()
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}

	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialPageLayout, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/page/layout_index`)
	ctx.Set(`title`, ctx.T(`编辑布局`))
	return ctx.Render(`official/page/layout_edit`, common.Err(ctx, err))
}

func LayoutDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelPage.NewLayout(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/page/layout_index`))
}
