package group

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
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
	_, err := common.PagingWithListerCond(ctx, m, cond)
	ret := common.Err(ctx, err)
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
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonGroup, uint64(id))
			}
		}
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonGroup,
		formbuilder.ConfigFile(`official/customer/group/edit`),
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
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonGroup, uint64(m.Id), i18nm.OptionContentType(`description`, `text`))
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
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

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
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/group/index`))
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
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/group/index`))
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
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/group/index`))
}
