package wallet

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	registryWallet "github.com/coscms/webfront/registry/wallet"
)

func Recharge(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	rechargeCfg := config.Setting(`base`).Get(`recharge`).(*modelCustomer.WalletSettings)
	if rechargeCfg == nil {
		return ctx.NewError(code.Unsupported, `暂未开启充值功能`)
	}
	m := modelCustomer.NewWallet(ctx)
	amount := ctx.Formx(`amount`).Float64()
	if amount <= 0.0 {
		amount = rechargeCfg.DefaultAmount
	}
	if amount < rechargeCfg.MinAmount {
		amount = rechargeCfg.MinAmount
	}
	ctx.Request().Form().Set("amount", param.AsString(amount))
	money, err := m.GetBalance(`money`, customer.Id)
	if err != nil {
		return err
	}
	err = registryWallet.RechargePage.Fire(ctx)
	if err != nil {
		return err
	}
	ctx.Set(`money`, money)
	ctx.Set(`rechargeCfg`, rechargeCfg)
	ctx.Set(`activeURL`, `/user/wallet`)
	ctx.Set(`rechargePage`, registryWallet.RechargePage)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	return ctx.Render(`user/wallet/recharge`, common.Err(ctx, err))
}
