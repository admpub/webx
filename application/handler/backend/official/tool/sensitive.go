package tool

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webfront/library/sensitive"
	modelSensitive "github.com/coscms/webfront/model/official/sensitive"
)

func sensitiveIndex(ctx echo.Context) error {
	m := modelSensitive.NewSensitive(ctx)
	cond := db.NewCompounds()
	nsql.SelectPageCond(ctx, cond, `id`, `words%`)
	typ := ctx.Formx(`type`).String()
	if len(typ) > 0 {
		cond.AddKV(`type`, typ)
	}
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`official/tool/sensitive/index`, common.Err(ctx, err))
}

func sensitiveAdd(ctx echo.Context) error {
	var err error
	m := modelSensitive.NewSensitive(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonSensitive, echo.ExcludeFieldName(`updated`))
		if err == nil {
			var added []string
			added, err = common.BatchAdd(ctx, `words`, m, func(words *string) error {
				m.Id = 0
				if m.Type == `bad` && m.Disabled == `N` {
					sensitive.AddWord(*words)
				}
				return nil
			})
			if err == nil && len(added) == 0 {
				err = ctx.E(`关键词不能为空`)
			}
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				if m.Type == `noise` && m.Disabled == `N` {
					sensitive.Reset()
				}
				return ctx.Redirect(backend.URLFor(`/tool/sensitive/index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonSensitive, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/tool/sensitive/index`)
	ctx.Set(`title`, ctx.T(`添加敏感词`))
	return ctx.Render(`official/tool/sensitive/edit`, common.Err(ctx, err))
}

func sensitiveEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := modelSensitive.NewSensitive(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonSensitive, echo.ExcludeFieldName(`created`))
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				sensitive.Reset()
				return ctx.Redirect(backend.URLFor(`/tool/sensitive/index`))
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
			if disabled == `Y` {
				sensitive.DelWord(m.Words)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}

	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonSensitive, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/tool/sensitive/index`)
	ctx.Set(`title`, ctx.T(`修改敏感词`))
	return ctx.Render(`official/tool/sensitive/edit`, common.Err(ctx, err))
}

func sensitiveDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelSensitive.NewSensitive(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/tool/sensitive/index`))
}
