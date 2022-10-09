package shorturl

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	modelShorturl "github.com/admpub/webx/application/model/official/shorturl"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

// DomainIndex 域名管理
func DomainIndex(ctx echo.Context) error {
	var err error
	m := modelShorturl.NewShortURL(ctx)
	cond := db.Compounds{}
	domain := ctx.Formx(`q`).String()
	if len(domain) > 0 {
		cond.AddKV(`domain`, domain)
	}
	sorts := common.Sorts(ctx, `official_short_url_domain`, `-id`)
	_, err = common.NewLister(m.Domain, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	ctx.Set(`listData`, m.Domain.Objects())
	return ctx.Render(`official/short_url/domain_index`, common.Err(ctx, err))
}

// DomainEdit 修改短网址
func DomainEdit(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.E(`参数“%s”值无效`, `id`)
	}
	m := modelShorturl.NewShortURL(ctx)
	err := m.Domain.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`域名不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		if err = ctx.MustBind(m.Domain, echo.IncludeFieldName(`disabled`)); err != nil {
			goto END
		}
		err = m.Domain.UpdateField(nil, `disabled`, m.Domain.Disabled, `id`, m.Domain.Id)
		//err = m.Edit(nil, `id`, m.Domain.Id)
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(handler.URLFor(`/official/short_url/domain_index`))
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Domain.Disabled = disabled
			data := ctx.Data()
			err = m.Domain.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.Domain, ``, echo.LowerCaseFirstLetter)
	}

END:
	ctx.Set(`activeURL`, `/official/short_url/domain_index`)
	ctx.Set(`title`, ctx.T(`修改域名`))
	return ctx.Render(`official/short_url/domain_edit`, common.Err(ctx, err))
}

// DomainDelete 删除域名
func DomainDelete(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.E(`参数“%s”值无效`, `id`)
	}
	m := modelShorturl.NewShortURL(ctx)
	err := m.Domain.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`域名不存在`)
		}
		return err
	}
	err = m.Domain.Delete(nil, `id`, id)
	if err != nil {
		return err
	}
	common.SendOk(ctx, ctx.T(`删除成功`))
	return ctx.Redirect(handler.URLFor(`/official/short_url/domain_index`))
}
