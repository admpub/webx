package tool

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
)

func AreaCountryIndex(ctx echo.Context) error {
	m := official.NewAreaCountry(ctx)
	cond := db.NewCompounds()
	nsql.SelectPageCond(ctx, cond, `id`, `name%,abbr%`)
	err := m.ListPage(cond, `sort`)
	if err != nil {
		return err
	}
	list := m.Objects()
	if ctx.Form(`select2`) == `1` {
		data := ctx.Data()
		ctx.Set(`listData`, list)
		r := make([]echo.H, len(list))
		for k, v := range list {
			r[k] = echo.H{`id`: v.Abbr, `text`: v.Name}
		}
		ctx.Set(`listData`, r)
		return ctx.JSON(data.SetData(ctx.Stored))
	}
	ctx.Set(`listData`, list)

	ctx.Set(`title`, ctx.T(`国家管理`))
	ctx.Set(`activeURL`, `/tool/area/index`)
	return ctx.Render(`official/tool/area/country_index`, common.Err(ctx, err))
}

func AreaCountryAdd(ctx echo.Context) error {
	var err error
	m := official.NewAreaCountry(ctx)
	form := formbuilder.New(ctx,
		m.OfficialCommonAreaCountry,
		formbuilder.ConfigFile(`official/tool/area/country_edit`),
		formbuilder.AllowedNames(
			`countryAbbr`, `name`, `abbr`, `areaIds`, `sort`,
		),
	)
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonAreaCountry, uint64(m.Id))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/tool/area/group_index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`添加地区分组`))
	return ctx.Render(`official/tool/area/group_edit`, err)
}

func AreaCountryEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, `参数错误`).SetZone(`id`)
	}
	m := official.NewAreaCountry(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsGet() {
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonAreaCountry, uint64(id))
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonAreaCountry,
		formbuilder.ConfigFile(`official/tool/area/group_edit`),
		formbuilder.AllowedNames(
			`countryAbbr`, `name`, `abbr`, `areaIds`, `sort`,
		),
	)
	form.OnPost(func() error {
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonAreaCountry, uint64(m.Id))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/tool/area/group_index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`编辑地区分组`))
	return ctx.Render(`official/tool/area/group_edit`, err)
}

func AreaCountryDelete(ctx echo.Context) error {
	id := ctx.FormxValues(`id`).Unique().Uint(param.IsGreaterThanZeroElement)
	if len(id) == 0 {
		return ctx.NewError(code.InvalidParameter, `请选择要删除的项`).SetZone(`id`)
	}
	m := official.NewAreaCountry(ctx)
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
