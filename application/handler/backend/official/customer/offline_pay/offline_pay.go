package group_package

import (
	"slices"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/offlinepay"
	"github.com/coscms/webfront/model/author"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

type OfflinePayWithCustomer struct {
	*dbschema.OfficialCustomerOfflinePay
	Customer *author.Customer `db:"-,relation=id:customer_id|gtZero,columns=id&name&avatar" json:",omitempty"`
}

func Index(ctx echo.Context) error {
	m := modelCustomer.NewOfflinePay(ctx)
	cond := db.NewCompounds()
	customerID := ctx.Formx(`customerId`).Uint64()
	if customerID > 0 {
		cond.AddKV(`customer_id`, customerID)
	}

	if targetType := ctx.Form(`targetType`); len(targetType) > 0 {
		cond.AddKV(`target_type`, targetType)
	}

	if targetID := ctx.Formx(`targetId`).Uint64(); targetID > 0 {
		cond.AddKV(`target_id`, targetID)
	}

	if payMethod := ctx.Form(`payMethod`); len(payMethod) > 0 {
		cond.AddKV(`pay_method`, payMethod)
	}

	if status := ctx.Form(`status`); len(status) > 0 {
		cond.AddKV(`status`, status)
	}

	list := []*OfflinePayWithCustomer{}
	err := m.ListPageByOffsetAs(&list, cond, `-id`)
	ret := common.Err(ctx, err)
	ctx.Set(`listData`, list)
	ctx.Set(`targetTypes`, modelCustomer.OfflinePayTargetTypes.Slice())
	statusList := modelCustomer.OfflinePayStatuses.Slice()
	ctx.Set(`statusList`, statusList)
	ctx.SetFunc(`targetTypeName`, modelCustomer.OfflinePayTargetTypes.Get)
	ctx.SetFunc(`ownershipInfo`, func(targetType string, ownershipID uint64) modelCustomer.OwnershipInfo {
		item := modelCustomer.OfflinePayTargetTypes.GetItem(targetType)
		if item == nil || item.X == nil {
			return modelCustomer.OwnershipInfo{}
		}
		return item.X.OwnershipInfo(ctx, ownershipID)
	})
	ctx.SetFunc(`statusName`, modelCustomer.OfflinePayStatuses.Get)
	payMethods := offlinepay.GetMethods(nil)
	ctx.Set(`payMethods`, offlinepay.GetMethods(nil))
	ctx.Set(`statusSource`, echo.KVListToEditableSource(ctx, echo.KVList(statusList)))
	ctx.SetFunc(`payMethodName`, func(v string) string {
		for _, item := range payMethods {
			if item.K == v {
				return item.V
			}
		}
		return ``
	})
	ctx.Set(`title`, ctx.T(`线下转账列表`))
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
	ctx.Set(`payMethods`, offlinepay.GetMethods(nil))
	ctx.Set(`targetTypes`, modelCustomer.OfflinePayTargetTypes.Slice())
	ctx.Set(`statusList`, modelCustomer.OfflinePayStatuses.Slice())
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
		m.Id = id
		err = m.Edit(nil, `id`, id)
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
			case modelCustomer.OfflinePayStatusPending:
				err = m.SetPending()
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
	ctx.Set(`payMethods`, offlinepay.GetMethods(nil))
	ctx.Set(`targetTypes`, modelCustomer.OfflinePayTargetTypes.Slice())
	ctx.Set(`statusList`, modelCustomer.OfflinePayStatuses.Slice())
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
