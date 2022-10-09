package wallet

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/webx/application/dbschema"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
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
	_, err = handler.NewLister(flowM, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, `customer_id`, customer.Id).Paging(ctx)
	ctx.Set(`list`, flowM.Objects())
	ctx.Set(`activeURL`, `/user/wallet`)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	return ctx.Render(`user/wallet/index`, handler.Err(ctx, err))
}
