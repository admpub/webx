package advert

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	modelAdvert "github.com/admpub/webx/application/model/official/advert"
)

func PositionIndex(ctx echo.Context) error {
	q := ctx.Form(`q`)
	contype := ctx.Form(`contype`)
	m := modelAdvert.NewAdPosition(ctx)
	cond := []db.Compound{}
	if len(contype) > 0 {
		cond = append(cond, db.Cond{`contype`: contype})
	}
	if len(q) > 0 {
		cond = append(cond, db.Cond{`name`: db.Like(`%` + q + `%`)})
	}
	ctx.Request().Form().Set(`pageSize`, `10`)
	list := []*modelAdvert.PositionWithRendered{}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(cond...)))
	for _, row := range list {
		row.Rendered = modelAdvert.Render(row)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`official/advert/position_index`, handler.Err(ctx, err))
}

func PositionAdd(ctx echo.Context) error {
	m := modelAdvert.NewAdPosition(ctx)
	var err error
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialAdPosition, echo.ExcludeFieldName(`updated`, `created`))
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/official/advert/position_index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialAdPosition, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/advert/position_index`)
	ctx.Set(`isEdit`, false)
	setPositionFormData(ctx)
	return ctx.Render(`official/advert/position_edit`, handler.Err(ctx, err))
}

func setPositionFormData(ctx echo.Context) {
	ctx.Set(`contypes`, modelAdvert.Contype.Slice())
}

func PositionEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelAdvert.NewAdPosition(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialAdPosition, echo.ExcludeFieldName(`updated`, `created`))
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/advert/position_index`))
			}
		}
	} else if ctx.IsAjax() {
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
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialAdPosition, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/official/advert/position_index`)
	ctx.Set(`isEdit`, true)
	setPositionFormData(ctx)
	return ctx.Render(`official/advert/position_edit`, handler.Err(ctx, err))
}

func PositionDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelAdvert.NewAdPosition(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/advert/position_index`))
}
