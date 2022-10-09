package tool

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
)

func AreaGroupIndex(ctx echo.Context) error {
	m := official.NewAreaGroup(ctx)
	cond := db.NewCompounds()
	common.SelectPageCond(ctx, cond, `id`, `name%,abbr%`)
	var err error
	if ctx.Form(`select2`) == `1` {
		var list []*dbschema.OfficialCommonAreaGroup
		data := ctx.Data()
		list, err = m.ListPage(cond, `sort`, `id`)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		ctx.Set(`listData`, list)
		r := make([]echo.H, len(list))
		for k, v := range list {
			r[k] = echo.H{`id`: v.Id, `text`: v.Name}
		}
		ctx.Set(`listData`, r)
		return ctx.JSON(data.SetData(ctx.Stored))
	}
	var list []*official.AreaGroupExt
	list, err = m.ListPageWithExt(cond, `sort`, `id`)
	ctx.Set(`listData`, list)

	ctx.Set(`title`, ctx.T(`地区分组`))
	ctx.Set(`activeURL`, `/tool/area/index`)
	return ctx.Render(`official/tool/area/group_index`, handler.Err(ctx, err))
}

func AreaGroupAdd(ctx echo.Context) error {
	var err error
	m := official.NewAreaGroup(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonAreaGroup, echo.ExcludeFieldName(`updated`))
		if err == nil {
			var added []string
			added, err = common.BatchAdd(ctx, `name`, m, func(_ *string) error {
				m.Id = 0
				return nil
			})
			if err == nil && len(added) == 0 {
				err = ctx.E(`地区名称不能为空`)
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/tool/area/group_index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonAreaGroup, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`添加地区分组`))
	return ctx.Render(`official/tool/area/group_edit`, err)
}

func AreaGroupEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := official.NewAreaGroup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonAreaGroup, echo.ExcludeFieldName(`created`, `updated`))
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/tool/area/group_index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonAreaGroup, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`编辑地区分组`))
	return ctx.Render(`official/tool/area/group_edit`, err)
}

func AreaGroupDelete(ctx echo.Context) error {
	id := ctx.FormxValues(`id`).Uint(func(index int, value uint) bool {
		return value > 0
	})
	m := official.NewAreaGroup(ctx)
	var err error
	for _, _v := range id {
		if err = m.Delete(nil, db.Cond{`id`: _v}); err != nil {
			break
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/tool/area/group_index`))
}
