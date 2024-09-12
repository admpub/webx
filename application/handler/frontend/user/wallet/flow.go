package wallet

import (
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
	}
	pagination.SetPageDefaultSize(ctx, 20)
	_, err := common.NewLister(flowM, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()).Paging(ctx)
	ctx.Set(`list`, flowM.Objects())
	ctx.Set(`activeURL`, `/user/wallet`)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	return ctx.Render(`user/wallet/flow`, common.Err(ctx, err))
}
