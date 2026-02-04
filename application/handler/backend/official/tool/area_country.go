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
	err := m.ListPage(cond, `sort`, `id`)
	if err != nil {
		return err
	}
	list := m.Objects()
	if ctx.Form(`select2`) == `1` {
		data := ctx.Data()
		ctx.Set(`listData`, list)
		r := make([]echo.H, len(list))
		for k, v := range list {
			r[k] = echo.H{`id`: v.Id, `text`: v.Name, `abbr`: v.Abbr}
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
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonAreaCountry, uint64(id))
			}
		}
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonAreaCountry,
		formbuilder.ConfigFile(`official/tool/area/country_edit`),
		formbuilder.AllowedNames(
			`short`, `name`, `abbr`, `code`, `lng`, `lat`, `disabled`, `sort`,
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
		return ctx.Redirect(backend.URLFor(`/tool/area/country_index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)
	shortField := form.MultilingualField(config.FromFile().Language.Default, `short`, `short`)
	shortField.AddTag(`required`)

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`添加国家`))
	return ctx.Render(`official/tool/area/country_edit`, common.Err(ctx, err))
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
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonAreaCountry, uint64(id))
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonAreaCountry,
		formbuilder.ConfigFile(`official/tool/area/country_edit`),
		formbuilder.AllowedNames(
			`short`, `name`, `abbr`, `code`, `disabled`, `lng`, `lat`, `sort`,
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
		return ctx.Redirect(backend.URLFor(`/tool/area/country_index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)
	shortField := form.MultilingualField(config.FromFile().Language.Default, `short`, `short`)
	shortField.AddTag(`required`)

	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`编辑国家`))
	return ctx.Render(`official/tool/area/country_edit`, common.Err(ctx, err))
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

	return ctx.Redirect(backend.URLFor(`/tool/area/country_index`))
}
