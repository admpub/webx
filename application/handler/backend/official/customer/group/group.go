package group

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/webx/application/model/official"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Index(ctx echo.Context) error {
	m := official.NewGroup(ctx)
	cond := db.Cond{}
	t := ctx.Form(`type`)
	if len(t) > 0 {
		cond[`type`] = t
	}
	_, err := handler.PagingWithListerCond(ctx, m, cond)
	ret := handler.Err(ctx, err)
	list := m.Objects()
	tg := make([]*official.GroupAndType, len(list))
	for k, u := range list {
		tg[k] = &official.GroupAndType{
			OfficialCommonGroup: u,
			Type:                &echo.KV{},
		}
		if len(u.Type) < 1 {
			continue
		}
		if typ := official.GroupTypes.GetItem(u.Type); typ != nil {
			tg[k].Type = typ
		}
	}

	ctx.Set(`listData`, tg)
	ctx.Set(`groupTypes`, official.GroupTypes.Slice())
	ctx.Set(`type`, t)
	return ctx.Render(`official/customer/group/index`, ret)
}

func Add(ctx echo.Context) error {
	var err error
	m := official.NewGroup(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`用户组名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`用户组名称已经存在`)
		} else {
			err = ctx.MustBind(m.OfficialCommonGroup)
		}
		if err == nil {
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/customer/group/index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonGroup, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/customer/group/index`)
	ctx.Set(`groupTypes`, official.GroupTypes.Slice())
	return ctx.Render(`official/customer/group/edit`, err)
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := official.NewGroup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`用户组名称不能为空`)
		} else if y, e := m.ExistsOther(name, id); e != nil {
			err = e
		} else if y {
			err = ctx.E(`用户组名称已经存在`)
		} else {
			err = ctx.MustBind(m.OfficialCommonGroup, echo.ExcludeFieldName(`created`))
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/customer/group/index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonGroup, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/official/customer/group/index`)
	ctx.Set(`groupTypes`, official.GroupTypes.Slice())
	return ctx.Render(`official/customer/group/edit`, err)
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := official.NewGroup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/customer/group/index`))
}
