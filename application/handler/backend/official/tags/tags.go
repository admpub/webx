package tags

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
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
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonTags, FormFilter())
		if err == nil {
			m.Group = ctx.Form(`newGroup`)
			_, err = m.Add()
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/tags/index`))
			}
		}
	}

	ctx.Set(`activeURL`, `/official/tags/index`)
	ctx.Set(`isEdit`, false)
	return ctx.Render(`official/tags/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	var err error
	name := ctx.Form(`name`)
	group := ctx.Form(`group`)
	m := official.NewTags(ctx)
	err = m.Get(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`group`: group},
	))
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonTags, FormFilter(formfilter.Exclude(`name`)))
		if err == nil {
			m.Group = ctx.Form(`newGroup`)
			m.Name = name
			err = m.Edit(nil, db.Cond{`name`: name})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/tags/index`))
			}
		}
	} else if ctx.IsAjax() {
		display := ctx.Query(`display`)
		if len(display) > 0 {
			m.Display = display
			data := ctx.Data()
			err = m.UpdateField(nil, `display`, display, db.Cond{`name`: name})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonTags, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/tags/index`)
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
