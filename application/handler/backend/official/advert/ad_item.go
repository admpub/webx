package advert

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v4/application/handler"
	modelAdvert "github.com/admpub/webx/application/model/official/advert"
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
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()))
	for _, row := range list {
		row.Rendered = modelAdvert.Render(row)
	}
	ctx.Set(`listData`, list)
	ctx.Set(`contypes`, modelAdvert.Contype.Slice())
	return ctx.Render(`official/advert/index`, handler.Err(ctx, err))
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
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialAdItem, formFilter())
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/official/advert/index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialAdItem, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/advert/index`)
	ctx.Set(`isEdit`, false)
	setFormData(ctx)
	return ctx.Render(`official/advert/edit`, handler.Err(ctx, err))
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
	m := modelAdvert.NewAdItem(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialAdItem, formFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/advert/index`))
			}
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
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
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialAdItem, ``, echo.LowerCaseFirstLetter)
		var startDate, endDate string
		if m.OfficialAdItem.Start > 0 {
			startDate = time.Unix(int64(m.OfficialAdItem.Start), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`start`, startDate)
		if m.OfficialAdItem.End > 0 {
			endDate = time.Unix(int64(m.OfficialAdItem.End), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`end`, endDate)
	}

	ctx.Set(`activeURL`, `/official/advert/index`)
	ctx.Set(`isEdit`, true)
	setFormData(ctx)
	return ctx.Render(`official/advert/edit`, handler.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelAdvert.NewAdItem(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/advert/index`))
}
