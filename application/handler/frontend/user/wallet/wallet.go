package wallet

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func Index(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	wallet := modelCustomer.NewWallet(ctx)
	assets, err := wallet.ListCustomerAllAssets(customer.Id)
	if err != nil {
		return err
	}
	ctx.Set(`assets`, assets)
	ctx.Set(`assetTypes`, modelCustomer.AssetTypeList())
	pagination.SetPageDefaultSize(ctx, 20)
	flowM := dbschema.NewOfficialCustomerWalletFlow(ctx)
	_, err = common.NewLister(flowM, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, `customer_id`, customer.Id).Paging(ctx)
	ctx.Set(`list`, flowM.Objects())
	ctx.Set(`activeURL`, `/user/wallet`)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	ctx.SetFunc(`formatAnyAssetAmount`, modelCustomer.MakeAnyAssetAmountFormatter(ctx))
	return ctx.Render(`user/wallet/index`, common.Err(ctx, err))
}
