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

func AreaIndex(ctx echo.Context) error {
	m := official.NewArea(ctx)
	cond := db.NewCompounds()
	pid := ctx.Formx(`pid`).Uint()
	cond.AddKV(`pid`, pid)
	countryAbbr := ctx.Form(`countryAbbr`)
	if len(countryAbbr) > 0 {
		cond.Add(db.Cond{`country_abbr`: countryAbbr})
	}
	nsql.SelectPageCond(ctx, cond, `id`, `name%,pinyin%`)
	queryMW := func(r db.Result) db.Result {
		return r.OrderBy(`id`)
	}
	_, err := common.NewLister(m.OfficialCommonArea, nil, queryMW, cond.And()).Paging(ctx)
	ret := common.Err(ctx, err)
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
	var pids []uint
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonArea, uint64(id))
				pids, _ = m.PositionIds(m.Pid)
			}
		}

	}
	form := formbuilder.New(ctx,
		m.OfficialCommonArea,
		formbuilder.ConfigFile(`official/tool/area/edit`),
		formbuilder.AllowedNames(
			`countryAbbr`, `name`, `short`, `merged`, `pid`, `pinyin`, `code`, `zip`, `lng`, `lat`,
		),
	)
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonArea, uint64(m.Id))
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

	ctx.Set(`pids`, pids)
	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`添加地区`))
	return ctx.Render(`official/tool/area/edit`, err)
}

func AreaEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, `参数错误`).SetZone(`id`)
	}
	m := official.NewArea(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	var pids []uint
	if ctx.IsGet() {
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonArea, uint64(id))
		pids, _ = m.PositionIds(m.Pid)
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonArea,
		formbuilder.ConfigFile(`official/tool/area/edit`),
		formbuilder.AllowedNames(
			`countryAbbr`, `name`, `short`, `merged`, `pid`, `pinyin`, `code`, `zip`, `lng`, `lat`,
		),
	)
	form.OnPost(func() error {
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonArea, uint64(m.Id))
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

	ctx.Set(`pids`, pids)
	ctx.Set(`activeURL`, `/tool/area/index`)
	ctx.Set(`title`, ctx.T(`编辑地区`))
	return ctx.Render(`official/tool/area/edit`, err)
}

func AreaDelete(ctx echo.Context) error {
	id := ctx.FormxValues(`id`).Uint(param.IsGreaterThanZeroElement)
	if len(id) == 0 {
		return ctx.NewError(code.InvalidParameter, `请选择要删除的项`).SetZone(`id`)
	}
	m := official.NewArea(ctx)
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

	return ctx.Redirect(backend.URLFor(`/tool/area/index`))
}
