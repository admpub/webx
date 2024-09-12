package agent

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	modelAgent "github.com/coscms/webfront/model/official/agent"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func productFormFilter() echo.FormDataFilter {
	return formfilter.Build(
		formfilter.Exclude(`sold`, `performance`, `created`, `updated`),
		formfilter.DateToTimestamp(`expired`),
	)
}

func ProductIndex(ctx echo.Context) error {
	if operation := ctx.Form(`operation`); operation == `selectProduct` {
		return selectPageProduct(ctx)
	}
	m := modelAgent.NewAgentProduct(ctx)
	cond := db.Cond{}
	list := []*modelAgent.AgentAndProduct{}
	_, err := common.PagingWithLister(ctx, common.NewLister(m, &list, func(r db.Result) db.Result {
		return r.Relation(`Agent`, modelCustomer.CusomterSafeFieldsSelector).OrderBy(`-id`)
	}, cond))
	if err == nil {
		err = modelAgent.WithProductInfo(ctx, list)
	}
	ctx.Set(`listData`, list)
	ctx.Set(`productTableList`, modelAgent.Source.Slice())
	ctx.SetFunc(`getProductTableName`, modelAgent.Source.Get)
	return ctx.Render(`official/agent/product_index`, common.Err(ctx, err))
}

func selectPageProduct(ctx echo.Context) error {
	productTable := ctx.Form(`productTable`)
	h := modelAgent.Source.GetSelectPageHandler(productTable)
	if h != nil {
		return h(ctx)
	}
	return nil
}

func ProductAdd(ctx echo.Context) error {
	if operation := ctx.Form(`operation`); operation == `selectProduct` {
		return selectPageProduct(ctx)
	}
	var err error
	m := modelAgent.NewAgentProduct(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerAgentProduct, productFormFilter())
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/agent/product_index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCustomerAgentProduct, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
				if m.Expired > 0 {
					ctx.Request().Form().Set(`expired`, time.Unix(int64(m.Expired), 0).Format(`2006-01-02`))
				}
			}
		}
	}

	ctx.Set(`activeURL`, `/official/agent/product_index`)
	ctx.Set(`productTableList`, modelAgent.Source.Slice())
	ctx.Set(`title`, ctx.T(`添加代理产品`))
	return ctx.Render(`official/agent/product_edit`, common.Err(ctx, err))
}

func ProductEdit(ctx echo.Context) error {
	if operation := ctx.Form(`operation`); operation == `selectProduct` {
		return selectPageProduct(ctx)
	}
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelAgent.NewAgentProduct(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerAgentProduct, productFormFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/agent/product_index`))
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
		echo.StructToForm(ctx, m.OfficialCustomerAgentProduct, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
		if m.Expired > 0 {
			ctx.Request().Form().Set(`expired`, time.Unix(int64(m.Expired), 0).Format(`2006-01-02`))
		}
	}

	ctx.Set(`activeURL`, `/official/agent/product_index`)
	ctx.Set(`productTableList`, modelAgent.Source.Slice())
	ctx.Set(`title`, ctx.T(`修改代理产品`))
	return ctx.Render(`official/agent/product_edit`, common.Err(ctx, err))
}

func ProductDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelAgent.NewAgentProduct(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/agent/product_index`))
}
