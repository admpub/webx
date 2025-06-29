package wallet

import (
	"html/template"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
)

func Flow(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	flowM := dbschema.NewOfficialCustomerWalletFlow(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`customer_id`, customer.Id)
	amountType := ctx.Form(`amountType`)
	assetType := ctx.Form(`assetType`)
	if len(amountType) > 0 {
		cond.AddKV(`amount_type`, amountType)
	}
	if len(assetType) > 0 {
		cond.AddKV(`asset_type`, assetType)
		formatter := modelCustomer.MakeAssetAmountFormatter(ctx, assetType)
		ctx.SetFunc(`formatAnyAssetAmount`, func(assetType string, amount float64) template.HTML {
			return formatter(amount)
		})
	} else {
		ctx.SetFunc(`formatAnyAssetAmount`, modelCustomer.MakeAnyAssetAmountFormatter(ctx))
	}
	pagination.SetPageDefaultSize(ctx, 20)
	err := flowM.ListPageByOffset(cond, `-id`)
	ctx.Set(`list`, flowM.Objects())
	ctx.Set(`activeURL`, `/user/wallet`)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	return ctx.Render(`user/wallet/flow`, common.Err(ctx, err))
}
