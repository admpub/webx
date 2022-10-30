package wallet

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func FlowIndex(ctx echo.Context) error {
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
	sourceType := ctx.Form(`sourceType`)
	if len(sourceType) > 0 {
		cond.AddKV(`source_type`, sourceType)
	}
	sourceTable := ctx.Form(`sourceTable`)
	if len(sourceTable) > 0 {
		cond.AddKV(`source_table`, sourceTable)
	}
	sourceID := ctx.Form(`sourceId`)
	if len(sourceID) > 0 {
		cond.AddKV(`source_id`, sourceID)
	}
	sourceCustomer := ctx.Formx(`sourceCustomer`).Uint64()
	if sourceCustomer > 0 {
		cond.AddKV(`source_customer`, sourceCustomer)
	}
	tradeNo := ctx.Form(`tradeNo`)
	if len(tradeNo) > 0 {
		cond.AddKV(`trade_no`, tradeNo)
	}
	typ := ctx.Form(`type`)
	if len(typ) > 0 {
		if typ == `income` { // 收入
			cond.AddKV(`amount`, db.Gte(0))
		} else { // 消费
			cond.AddKV(`amount`, db.Lt(0))
		}
	}
	amountType := ctx.Form(`amountType`)
	if len(amountType) > 0 {
		cond.AddKV(`amount_type`, amountType)
	}
	list, err := m.FlowListPage(cond, `-id`)
	ctx.Set(`listData`, list)
	assetTypes := modelCustomer.AssetTypeList()
	ctx.Set(`assetTypes`, assetTypes)
	ctx.Set(`assetType`, assetType)
	ctx.Set(`customer`, cs.ClearPasswordData())
	ctx.Set(`activeURL`, `/official/customer/index`)
	return ctx.Render(`official/customer/wallet/flow`, handler.Err(ctx, err))
}
