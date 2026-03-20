package group_package

import (
	"slices"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func Index(ctx echo.Context) error {
	targetType := ctx.Form(`targetType`)
	targetID := ctx.Formx(`targetId`).Uint64()
	m := modelCustomer.NewOfflinePay(ctx)
	cond := db.NewCompounds()
	if len(targetType) > 0 {
		cond.AddKV(`target_type`, targetType)
	}
	if targetID > 0 {
		cond.AddKV(`target_id`, targetID)
	}
	customerID := ctx.Formx(`customerId`).Uint64()
	if customerID > 0 {
		cond.AddKV(`customer_id`, customerID)
	}
	err := m.ListPageByOffset(cond, `-id`)
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`official/customer/offline_pay/index`, ret)
}

func Add(ctx echo.Context) error {
	var err error
	var id uint64
	m := modelCustomer.NewOfflinePay(ctx)
	if ctx.IsPost() {
		err = ctx.MustBindAndValidate(m.OfficialCustomerOfflinePay, echo.ExcludeFieldName(`id`, `status`, `created`, `updated`))
		if err != nil {
			goto END
		}

		_, err = m.Add()
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/customer/offline_pay/index`))
	}

	if id = ctx.Formx(`copyId`).Uint64(); id > 0 {
		err = m.Get(nil, `id`, id)
		if err != nil {
			return err
		}
		m.Id = 0
	}

END:
	ctx.Set(`activeURL`, `/official/customer/offline_pay/index`)
	ctx.Set(`title`, ctx.T(`添加线下转账`))
	return ctx.Render(`official/customer/offline_pay/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, `参数错误`).SetZone(`id`)
	}
	m := modelCustomer.NewOfflinePay(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBindAndValidate(m.OfficialCustomerOfflinePay, echo.ExcludeFieldName(`id`, `status`, `created`, `updated`))
		if err != nil {
			goto END
		}

		_, err = m.Add()
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/customer/offline_pay/index`))
	}
	if ctx.IsAjax() {

		status := ctx.Query(`status`)
		if len(status) > 0 {
			if !slices.Contains(modelCustomer.OfflinePayStatusAll, status) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`status`)
			}
			if m.Status == modelCustomer.OfflinePayStatusVerified {
				return ctx.NewError(code.DataStatusIncorrect, `数据已经核实过了，不能修改`).SetZone(`status`)
			}
			switch status {
			case modelCustomer.OfflinePayStatusVerified:
				err = m.SetVerified()
			case modelCustomer.OfflinePayStatusInvalid:
				err = m.SetInvalid()
			default:
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`status`)
			}
			data := ctx.Data()
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	}
	echo.StructToForm(ctx, m.OfficialCustomerOfflinePay, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`activeURL`, `/official/customer/offline_pay/index`)
	ctx.Set(`title`, ctx.T(`编辑线下转账`))
	return ctx.Render(`official/customer/offline_pay/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewOfflinePay(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/offline_pay/index`))
}
