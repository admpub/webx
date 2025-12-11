package advert

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/model/i18nm"
	modelAdvert "github.com/coscms/webfront/model/official/advert"
)

func Index(ctx echo.Context) error {
	q := ctx.Form(`q`)
	contype := ctx.Form(`contype`)
	positionID := ctx.Formx(`positionId`).Uint64()
	m := modelAdvert.NewAdItem(ctx)
	cond := db.NewCompounds()
	if len(contype) > 0 {
		cond.Add(db.Cond{`contype`: contype})
	}
	if len(q) > 0 {
		cond.Add(db.Cond{`name`: db.Like(`%` + q + `%`)})
	}
	ctx.Request().Form().Set(`pageSize`, `10`)
	list := []*modelAdvert.ItemAndPosition{}
	sorts := []interface{}{}
	if positionID > 0 {
		cond.Add(db.Cond{`position_id`: positionID})
		pm := modelAdvert.NewAdPosition(ctx)
		pm.Get(nil, db.Cond{`id`: positionID})
		ctx.Set(`positionInfo`, pm.OfficialAdPosition)
	} else {
		sorts = append(sorts, `-position_id`)
	}
	sorts = append(sorts, `sort`, `id`)
	_, err := common.PagingWithLister(ctx, common.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()))
	for _, row := range list {
		row.Rendered = modelAdvert.Render(row)
	}
	ctx.Set(`listData`, list)
	ctx.Set(`contypes`, modelAdvert.Contype.Slice())
	return ctx.Render(`official/advert/index`, common.Err(ctx, err))
}

func formFilter() echo.FormDataFilter {
	return formfilter.Build(
		formfilter.StartDateToTimestamp(`Start`),
		formfilter.EndDateToTimestamp(`End`),
		formfilter.Exclude(`updated`, `created`),
	)
}

func Add(ctx echo.Context) error {
	m := modelAdvert.NewAdItem(ctx)
	var err error
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialAdItem, id)
			} else {
				m.Sort = 500
			}
		} else {
			m.Sort = 500
		}
	}
	form := formbuilder.New(ctx,
		m.OfficialAdItem,
		formbuilder.ConfigFile(`official/advert/edit`),
		formbuilder.AllowedNames(
			`name`, `positionId`, `contype`, `mode`, `content`, `url`, `start`, `end`, `sort`, `disabled`,
		),
	)
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialAdItem, m.Id,
			i18nm.OptionContentType(`content`, m.Contype),
		)
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(backend.URLFor(`/official/advert/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()

	ctx.Set(`activeURL`, `/official/advert/index`)
	ctx.Set(`isEdit`, false)
	setFormData(ctx)
	return ctx.Render(`official/advert/edit`, common.Err(ctx, err))
}

func setFormData(ctx echo.Context) {
	ctx.Set(`contypes`, modelAdvert.Contype.Slice())
	ctx.Set(`modes`, modelAdvert.AdMode.Slice())
	pm := modelAdvert.NewAdPosition(ctx)
	pm.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`disabled`, `id`)
	}, 0, -1)
	ctx.Set(`positionList`, pm.Objects())
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, `参数错误`).SetZone(`id`)
	}
	m := modelAdvert.NewAdItem(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsGet() {
		if ctx.IsAjax() {
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
		}
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialAdItem, id)
	}
	form := formbuilder.New(ctx,
		m.OfficialAdItem,
		formbuilder.ConfigFile(`official/advert/edit`),
		formbuilder.AllowedNames(
			`name`, `positionId`, `contype`, `mode`, `content`, `url`, `start`, `end`, `sort`, `disabled`,
		),
	)
	form.OnPost(func() error {
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialAdItem, m.Id,
			i18nm.OptionContentType(`content`, m.Contype),
		)
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/advert/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()

	ctx.Set(`activeURL`, `/official/advert/index`)
	ctx.Set(`isEdit`, true)
	setFormData(ctx)
	return ctx.Render(`official/advert/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelAdvert.NewAdItem(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/advert/index`))
}
