package wallet

import (
	"github.com/coscms/webcore/library/common"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
)

// PrepaidCard 使用充值卡充值
func PrepaidCard(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	if err := ctx.Begin(); err != nil {
		return err
	}
	m := modelCustomer.NewWallet(ctx)
	cardM := modelCustomer.NewPrepaidCard(ctx)
	cardNumber := ctx.Formx(`cardNumber`).String()
	cardPassword := ctx.Formx(`cardPassword`).String()
	err := cardM.UseCard(customer.Id, cardNumber, cardPassword)
	if err != nil {
		ctx.Rollback()
		return err
	}
	m.Flow.CustomerId = customer.Id
	m.Flow.AssetType = modelCustomer.AssetTypeMoney
	m.Flow.AmountType = modelCustomer.AmountTypeBalance
	m.Flow.Amount = float64(cardM.Amount)
	m.Flow.SourceType = `recharge`
	m.Flow.SourceTable = `official_prepaid_card`
	m.Flow.SourceId = cardM.Id
	m.Flow.TradeNo = ``
	m.Flow.Status = modelCustomer.FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
	m.Flow.Description = `使用充值卡充值`
	err = m.AddRepeatableFlow()
	if err != nil {
		ctx.Rollback()
		return err
	}
	ctx.Commit()
	common.SendOk(ctx, ctx.T(`操作成功`))
	next := ctx.Form(`next`)
	if len(next) == 0 {
		next = `/user/wallet`
	}
	return ctx.Redirect(next)
}
