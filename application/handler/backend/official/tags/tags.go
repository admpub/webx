package tags

import (
	"github.com/coscms/forms/fields"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
)

func FormFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(options, formfilter.Exclude(`num`))
	return formfilter.Build(options...)
}

func Index(ctx echo.Context) error {
	name := ctx.Form(`name`, ctx.Form(`searchValue`))
	group := ctx.Form(`group`)
	m := official.NewTags(ctx)
	cond := db.NewCompounds()
	if len(group) > 0 {
		cond.AddKV(`group`, group)
	}
	if len(name) > 0 {
		cond.AddKV(`name`, db.Like(name+`%`))
	}
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r //.OrderBy(`-id`)
	}, cond.And()))
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`official/tags/index`, ret)
}

func Add(ctx echo.Context) error {
	var err error
	m := official.NewTags(ctx)
	form := formbuilder.New(ctx,
		m.OfficialCommonTags,
		formbuilder.ConfigFile(`official/tags/edit`),
		formbuilder.AllowedNames(`name`, `group`, `display`),
	)
	form.OnPost(func() error {
		m.Group = ctx.Form(`newGroup`)
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(m.OfficialCommonTags, uint64(m.Id))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(backend.URLFor(`/official/tags/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()

	ctx.Set(`activeURL`, `/official/tags/index`)
	ctx.Set(`groups`, official.TagGroups.Slice())
	ctx.Set(`isEdit`, false)
	return ctx.Render(`official/tags/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	var err error
	name := ctx.Form(`name`)
	group := ctx.Form(`group`)
	m := official.NewTags(ctx)
	cond := db.And(
		db.Cond{`name`: name},
		db.Cond{`group`: group},
	)
	err = m.Get(nil, cond)
	if err != nil {
		return err
	}
	if ctx.IsGet() {
		if ctx.IsAjax() {
			display := ctx.Query(`display`)
			if len(display) > 0 {
				m.Display = display
				data := ctx.Data()
				err = m.UpdateField(nil, `display`, display, cond)
				if err != nil {
					data.SetError(err)
					return ctx.JSON(data)
				}
				data.SetInfo(ctx.T(`操作成功`))
				return ctx.JSON(data)
			}
		}
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonTags,
		formbuilder.ConfigFile(`official/tags/edit`),
		formbuilder.AllowedNames(`group`, `display`),
		formbuilder.FormFilter(formfilter.Exclude(`name`)),
	)
	form.OnPost(func() error {
		m.Group = ctx.Form(`newGroup`)
		m.Name = name
		m.Display = common.GetBoolFlag(m.Display)
		err := m.UpdateFields(nil, echo.H{
			`group`:   m.Group,
			`display`: m.Display,
		}, cond)
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(m.OfficialCommonTags, uint64(m.Id))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/tags/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`).(*fields.Field)
	nameField.Type = `static`
	nameField.SetText(m.Name).SetTemplate(`static`)
	nameField.InitTemplate()
	echo.Dump(nameField)

	ctx.Set(`activeURL`, `/official/tags/index`)
	ctx.Set(`groups`, official.TagGroups.Slice())
	ctx.Set(`isEdit`, true)
	return ctx.Render(`official/tags/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	name := ctx.Form(`name`)
	group := ctx.Form(`group`)
	m := official.NewTags(ctx)
	err := m.Delete(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`group`: group},
	))
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/tags/index`))
}
