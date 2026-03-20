package wallet

import (
	"time"

	"github.com/admpub/dateparse"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/offlinepay"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
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

func (r RequestRechargeOffline) Apply(m *dbschema.OfficialCustomerOfflinePay) error {
	m.PayAccount = r.OfflinePayAccount
	m.PayAmount = r.OfflinePayAmount
	m.PayMethod = r.OfflinePayMethod
	m.PayBankBranch = r.OfflinePayBankBranch
	m.PayOwner = r.OfflinePayOwner
	m.PayTransactionNo = r.OfflinePayTransactionNo
	payTime, err := r.PayTime()
	if err != nil {
		return err
	}
	if !payTime.IsZero() {
		m.PayTime = uint(payTime.Unix())
	}
	return err
}

func (r RequestRechargeOffline) BeforeVadidate(ctx echo.Context) error {
	if offlinepay.GetMethod(r.OfflinePayMethod, nil) == nil {
		return ctx.NewError(code.InvalidParameter, `不支持的付款方式: %s`, r.OfflinePayMethod).SetZone(`offlinePayMethod`)
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
	m := modelCustomer.NewOfflinePay(ctx)
	m.CustomerId = customer.Id
	m.TargetType = modelCustomer.OfflinePayTargetTypeRecharge
	err := requestData.Apply(m.OfficialCustomerOfflinePay)
	if err != nil {
		return err
	}
	m.Status = modelCustomer.OfflinePayStatusPending
	_, err = m.Add()
	if err != nil {
		return err
	}
	common.SendOk(ctx, ctx.T(`操作成功`))
	next := ctx.Form(`next`)
	if len(next) == 0 {
		next = `/user/wallet/offline_pay_history`
	}
	return ctx.Redirect(next)
}

// RechargeOfflineHistory 线下转账充值历史
func RechargeOfflineHistory(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	m := modelCustomer.NewOfflinePay(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`customer_id`, customer.Id)
	pagination.SetDefaultSize(ctx, 20)
	err := m.ListPage(cond, `-id`)
	ctx.Set(`list`, m.Objects())
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
	ctx.SetFunc(`payMethodName`, func(v string) string {
		for _, item := range payMethods {
			if item.K == v {
				return item.V
			}
		}
		return ``
	})
	ctx.Set(`activeURL`, `/user/wallet`)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	return ctx.Render(`user/wallet/offline_pay_history`, common.Err(ctx, err))
}
