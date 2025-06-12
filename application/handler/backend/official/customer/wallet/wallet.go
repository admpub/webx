package wallet

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func Index(ctx echo.Context) error {
	m := modelCustomer.NewWallet(ctx)
	cond := db.NewCompounds()
	customerID := ctx.Formx(`customerId`).Uint64()
	cs := modelCustomer.NewCustomer(ctx)
	if customerID > 0 {
		cond.AddKV(`customer_id`, customerID)
		err := cs.Get(nil, `id`, customerID)
		if err != nil {
			if err == db.ErrNoMoreRows {
				return ctx.E(`客户信息不存在(ID: %v)`, customerID)
			}
			return err
		}
	}
	assetType := ctx.Form(`assetType`)
	if len(assetType) > 0 {
		cond.AddKV(`asset_type`, assetType)
	}
	list, err := m.ListPage(cond, `asset_type`)
	if err == nil &&
		customerID > 0 && len(assetType) > 0 &&
		len(list) == 0 && modelCustomer.AssetTypes.Has(assetType) {
		walletData := dbschema.NewOfficialCustomerWallet(ctx)
		walletData.CustomerId = customerID
		walletData.AssetType = assetType
		list = []*modelCustomer.WalletExt{
			{
				OfficialCustomerWallet: walletData,
				Customer: &modelCustomer.CustomerBase{
					Id:     walletData.CustomerId,
					Name:   cs.Name,
					Gender: cs.Gender,
					Avatar: cs.Avatar,
				},
				AssetTypeName: modelCustomer.AssetTypes.Get(walletData.AssetType),
			},
		}
	}
	ctx.Set(`listData`, list)
	assetTypes := modelCustomer.AssetTypeList()
	ctx.Set(`assetTypes`, assetTypes)
	ctx.Set(`assetType`, assetType)
	ctx.Set(`customer`, cs.ClearPasswordData())
	ctx.Set(`activeURL`, `/official/customer/index`)
	return ctx.Render(`official/customer/wallet/index`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	user := backend.User(ctx)
	var (
		err  error
		cond = db.NewCompounds()
	)
	customerID := ctx.Formx(`customerId`).Uint64()
	assetType := ctx.Form(`assetType`)
	if customerID == 0 {
		return ctx.E(`%v无效`, `customerId`)
	}
	m := modelCustomer.NewWallet(ctx)
	cs := modelCustomer.NewCustomer(ctx)
	err = cs.Get(nil, `id`, customerID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return ctx.E(`客户信息不存在(ID: %v)`, customerID)
		}
	}
	if len(assetType) == 0 {
		goto END
	}
	if !modelCustomer.AssetTypes.Has(assetType) {
		return ctx.E(`资产类型不存在: %v`, assetType)
	}
	cond.AddKV(`customer_id`, customerID)
	cond.AddKV(`asset_type`, assetType)
	err = m.Get(nil, cond.And())
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
	}
	if ctx.IsPost() {
		changeBalance := ctx.Formx(`changeBalance`).Float64()
		changeFreeze := ctx.Formx(`changeFreeze`).Float64()
		changeBalanceReason := ctx.Formx(`changeBalanceReason`).String()
		changeFreezeReason := ctx.Formx(`changeFreezeReason`).String()
		if changeBalance != 0 {
			m.Flow.Reset()
			m.Flow.CustomerId = cs.Id
			m.Flow.AssetType = assetType
			m.Flow.AmountType = modelCustomer.AmountTypeBalance
			m.Flow.Amount = changeBalance
			m.Flow.SourceType = ``
			m.Flow.SourceTable = `user`
			m.Flow.SourceId = uint64(user.Id)
			m.Flow.TradeNo = ``
			m.Flow.Status = modelCustomer.FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
			m.Flow.Description = changeBalanceReason
			if err = m.AddRepeatableFlow(); err != nil {
				goto END
			}
		}
		if changeFreeze != 0 {
			m.Flow.Reset()
			m.Flow.CustomerId = cs.Id
			m.Flow.AssetType = assetType
			m.Flow.AmountType = modelCustomer.AmountTypeFreeze
			m.Flow.Amount = changeFreeze
			m.Flow.SourceType = ``
			m.Flow.SourceTable = `user`
			m.Flow.SourceId = uint64(user.Id)
			m.Flow.TradeNo = ``
			m.Flow.Status = modelCustomer.FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
			m.Flow.Description = changeFreezeReason
			if err = m.AddRepeatableFlow(); err != nil {
				goto END
			}
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/customer/wallet/index?customerId=` + ctx.Query(`customerId`)))
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCustomerWallet, ``, echo.LowerCaseFirstLetter)
		ctx.Request().Form().Set(`customerId`, fmt.Sprint(customerID))
		ctx.Request().Form().Set(`assetType`, assetType)
	}

END:
	ctx.Set(`activeURL`, `/official/customer/index`)
	assetTypes := modelCustomer.AssetTypeList()
	ctx.Set(`assetTypes`, assetTypes)
	ctx.Set(`customer`, cs.ClearPasswordData())
	ctx.Set(`assetType`, assetType)
	ctx.Set(`assetTypeName`, modelCustomer.AssetTypes.Get(m.AssetType))
	return ctx.Render(`official/customer/wallet/edit`, err)
}
