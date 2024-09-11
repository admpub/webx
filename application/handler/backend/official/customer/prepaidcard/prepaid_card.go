package prepaidcard

import (
	"time"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"

	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func formFilter() echo.FormDataFilter {
	return formfilter.Build(
		formfilter.StartDateToTimestamp(`Start`),
		formfilter.EndDateToTimestamp(`End`),
	)
}

func Index(ctx echo.Context) error {
	m := modelCustomer.NewPrepaidCard(ctx)
	cond := db.Compounds{}
	q := ctx.Formx(`q`).String()
	if len(q) == 0 {
		q = ctx.Form(`number`)
	}
	if len(q) > 0 {
		cond.AddKV(`number`, q)
	}
	status := ctx.Form(`status`)
	if len(status) > 0 {
		switch status {
		case `used`:
			cond.AddKV(`used`, db.Gt(0))
		case `unused`:
			cond.AddKV(`used`, 0)
		}
	}
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`official/customer/prepaid_card/index`, common.Err(ctx, err))
}

func Add(ctx echo.Context) error {
	user := backend.User(ctx)
	var err error
	if ctx.IsPost() {
		m := modelCustomer.NewPrepaidCard(ctx)
		err = ctx.MustBind(m.OfficialCustomerPrepaidCard, formFilter())
		if err == nil {
			count := ctx.Formx(`count`).Int()
			err = m.BatchGenerate(user.Id, count, m.Amount, m.SalePrice, m.Start, m.End, m.BgImage)
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/customer/prepaid_card/index`))
		}
	}
	ctx.Set(`activeURL`, `/official/customer/prepaid_card/index`)
	ctx.Set(`title`, ctx.T(`添加充值卡`))
	ctx.Set(`isAdd`, true)
	return ctx.Render(`official/customer/prepaid_card/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewPrepaidCard(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		common.SendFail(ctx, err.Error())
		return ctx.Redirect(backend.URLFor(`/official/customer/prepaid_card/index`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerPrepaidCard, formFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(backend.URLFor(`/official/customer/prepaid_card/index`))
		}
	} else {
		echo.StructToForm(ctx, m.OfficialCustomerPrepaidCard, ``, echo.LowerCaseFirstLetter)
		var startDate, endDate string
		if m.OfficialCustomerPrepaidCard.Start > 0 {
			startDate = time.Unix(int64(m.OfficialCustomerPrepaidCard.Start), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`start`, startDate)
		if m.OfficialCustomerPrepaidCard.End > 0 {
			endDate = time.Unix(int64(m.OfficialCustomerPrepaidCard.End), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`end`, endDate)
	}

	ctx.Set(`activeURL`, `/official/customer/prepaid_card/index`)
	ctx.Set(`title`, ctx.T(`修改充值卡`))
	ctx.Set(`isAdd`, false)
	return ctx.Render(`official/customer/prepaid_card/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewPrepaidCard(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/prepaid_card/index`))
}
