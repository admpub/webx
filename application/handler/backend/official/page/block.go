package page

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	modelPage "github.com/admpub/webx/application/model/official/page"
)

func BlockIndex(ctx echo.Context) error {
	m := modelPage.NewBlock(ctx)
	cond := db.NewCompounds()
	common.SelectPageCond(ctx, cond, `id`, `name%`)
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`official/page/block_index`, ret)
}

func BlockAdd(ctx echo.Context) error {
	var err error
	m := modelPage.NewBlock(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialPageBlock, echo.ExcludeFieldName(`updated`, `created`))
		if err == nil {
			_, err = m.Insert()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/page/block_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialPageBlock, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/page/block_index`)
	ctx.Set(`title`, ctx.T(`添加区块`))
	return ctx.Render(`official/page/block_edit`, handler.Err(ctx, err))
}

func BlockEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := modelPage.NewBlock(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialPageBlock, echo.ExcludeFieldName(`updated`, `created`))
		if err == nil {
			m.Id = id
			err = m.Update(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/page/block_index`))
			}
		}
	} else if ctx.IsAjax() {

		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
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
		echo.StructToForm(ctx, m.OfficialPageBlock, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/page/block_index`)
	ctx.Set(`title`, ctx.T(`编辑区块`))
	return ctx.Render(`official/page/block_edit`, handler.Err(ctx, err))
}

func BlockDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelPage.NewBlock(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/page/block_index`))
}
