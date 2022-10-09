package tool

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/model/official"
)

func AreaIndex(ctx echo.Context) error {
	m := official.NewArea(ctx)
	cond := db.NewCompounds()
	pid := ctx.Formx(`pid`).Uint()
	cond.AddKV(`pid`, pid)
	common.SelectPageCond(ctx, cond, `id`, `name%,pinyin%`)
	queryMW := func(r db.Result) db.Result {
		return r.OrderBy(`id`)
	}
	_, err := handler.NewLister(m.OfficialCommonArea, nil, queryMW, cond.And()).Paging(ctx)
	ret := handler.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	if ctx.Form(`select2`) == `1` {
		r := make([]echo.H, len(list))
		for k, v := range list {
			r[k] = echo.H{`id`: v.Id, `text`: v.Name}
		}
		ctx.Set(`listData`, r)
	}
	positions, err := m.Positions(pid)
	if err != nil {
		return err
	}

	ctx.Set(`title`, ctx.T(`地区列表`))
	ctx.Set(`positions`, positions)
	return ctx.Render(`official/tool/area/index`, ret)
}

func AreaAdd(ctx echo.Context) error {
	var err error
	m := official.NewArea(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonArea, echo.ExcludeFieldName(`updated`))
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
			return ctx.Redirect(handler.URLFor(`/tool/area/index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonArea, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
				pids, _ := m.PositionIds(m.Pid)
				for _, pid := range pids {
					ctx.Request().Form().Add("pids[]", param.AsString(pid))
				}
			}
		}
	}

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`添加地区`))
	return ctx.Render(`official/tool/area/edit`, err)
}

func AreaEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := official.NewArea(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonArea, echo.ExcludeFieldName(`created`, `updated`))
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/tool/area/index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonArea, ``, echo.LowerCaseFirstLetter)
		pids, _ := m.PositionIds(m.Pid)
		for _, pid := range pids {
			ctx.Request().Form().Add("pids[]", param.AsString(pid))
		}
	}

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`编辑地区`))
	return ctx.Render(`official/tool/area/edit`, err)
}

func AreaDelete(ctx echo.Context) error {
	id := ctx.FormxValues(`id`).Uint(func(index int, value uint) bool {
		return value > 0
	})
	m := official.NewArea(ctx)
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

	return ctx.Redirect(handler.URLFor(`/tool/area/index`))
}
