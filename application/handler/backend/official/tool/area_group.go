package tool

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/official"
)

func AreaGroupIndex(ctx echo.Context) error {
	m := official.NewAreaGroup(ctx)
	cond := db.NewCompounds()
	nsql.SelectPageCond(ctx, cond, `id`, `name%,abbr%`)
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
	return ctx.Render(`official/tool/area/group_index`, common.Err(ctx, err))
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
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/tool/area/group_index`))
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
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/tool/area/group_index`))
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
	id := ctx.FormxValues(`id`).Unique().Uint(param.IsGreaterThanZeroElement)
	if len(id) == 0 {
		return ctx.NewError(code.InvalidParameter, `请选择要删除的项`).SetZone(`id`)
	}
	m := official.NewAreaGroup(ctx)
	var err error
	for _, _v := range id {
		if err = m.Delete(nil, db.Cond{`id`: _v}); err != nil {
			break
		}
	}
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/tool/area/group_index`))
}
