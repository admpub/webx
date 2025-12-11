package advert

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/model/i18nm"
	modelAdvert "github.com/coscms/webfront/model/official/advert"
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
	_, err := common.PagingWithLister(ctx, common.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(cond...)))
	for _, row := range list {
		row.Rendered = modelAdvert.Render(row)
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`official/advert/position_index`, common.Err(ctx, err))
}

func PositionAdd(ctx echo.Context) error {
	m := modelAdvert.NewAdPosition(ctx)
	var err error
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialAdPosition, id)
			}
		}
	}
	form := formbuilder.New(ctx,
		m.OfficialAdPosition,
		formbuilder.ConfigFile(`official/advert/position_edit`),
		formbuilder.AllowedNames(
			`ident`, `name`, `width`, `height`, `contype`, `content`, `url`, `disabled`,
		),
	)
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialAdPosition, m.Id,
			i18nm.OptionContentType(`content`, m.Contype),
		)
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(backend.URLFor(`/official/advert/position_index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()

	ctx.Set(`activeURL`, `/official/advert/position_index`)
	ctx.Set(`isEdit`, false)
	setPositionFormData(ctx)
	return ctx.Render(`official/advert/position_edit`, common.Err(ctx, err))
}

func setPositionFormData(ctx echo.Context) {
	ctx.Set(`contypes`, modelAdvert.Contype.Slice())
}

func PositionEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, `参数错误`).SetZone(`id`)
	}
	m := modelAdvert.NewAdPosition(ctx)
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
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialAdPosition, id)
	}
	form := formbuilder.New(ctx,
		m.OfficialAdPosition,
		formbuilder.ConfigFile(`official/advert/position_edit`),
		formbuilder.AllowedNames(
			`ident`, `name`, `width`, `height`, `contype`, `content`, `url`, `disabled`,
		),
	)
	form.OnPost(func() error {
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialAdPosition, m.Id,
			i18nm.OptionContentType(`content`, m.Contype),
		)
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/advert/position_index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()

	ctx.Set(`activeURL`, `/official/advert/position_index`)
	ctx.Set(`isEdit`, true)
	setPositionFormData(ctx)
	return ctx.Render(`official/advert/position_edit`, common.Err(ctx, err))
}

func PositionDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelAdvert.NewAdPosition(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/advert/position_index`))
}
