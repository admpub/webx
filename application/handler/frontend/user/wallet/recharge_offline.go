package wallet

import (
	"time"

	"github.com/admpub/dateparse"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/library/offlinepay"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
)

type RequestRechargeOffline struct {
	OfflinePayMethod        string  `validate:"required"`
	OfflinePayAccount       string  `validate:"required"`
	OfflinePayAmount        float64 `validate:"required,min=0.01"`
	OfflinePayBankBranch    string  // 银行支行（线下银行转账时有效）
	OfflinePayTransactionNo string  // 交易订单号（线上转账时有效）
	OfflinePayTime          string  // 付款时间（可选）
	OfflinePayOwner         string  `validate:"required"`
}

func (r RequestRechargeOffline) BeforeVadidate(ctx echo.Context) error {
	if offlinepay.GetMethod(r.OfflinePayMethod, nil) {

	}
	return nil
}

func (r RequestRechargeOffline) PayTime() (time.Time, error) {
	if len(r.OfflinePayTime) == 0 {
		return time.Time{}, nil
	}
	return dateparse.ParseAny(r.OfflinePayTime)
}

// RechargeOffline 使用线下转账充值
func RechargeOffline(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	requestData := &RequestRechargeOffline{}
	if err := ctx.MustBindAndValidate(requestData); err != nil {
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
