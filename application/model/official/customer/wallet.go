package customer

import (
	"github.com/admpub/decimal"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/xdatabase"
	"github.com/admpub/webx/application/middleware/sessdata"
)

func NewWallet(ctx echo.Context) *Wallet {
	m := &Wallet{
		OfficialCustomerWallet: dbschema.NewOfficialCustomerWallet(ctx),
		Flow:                   dbschema.NewOfficialCustomerWalletFlow(ctx),
	}
	return m
}

type Wallet struct {
	*dbschema.OfficialCustomerWallet
	Flow *dbschema.OfficialCustomerWalletFlow
}

func (f *Wallet) Add() (pk interface{}, err error) {
	return f.OfficialCustomerWallet.Insert()
}

func (f *Wallet) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCustomerWallet.Update(mw, args...)
}

func (f *Wallet) GetBalance(assetType string, customerID uint64) (float64, error) {
	return f.GetAssetBalance(assetType, customerID)
}

// AddRepeatableFlow 添加允许重复的流水记录
// unique_key: customer_id, asset_type, amount_type, source_type, source_table, source_id, number
func (f *Wallet) AddRepeatableFlow(flows ...*dbschema.OfficialCustomerWalletFlow) (err error) {
	flow := f.Flow
	if len(flows) > 0 {
		flow = flows[0]
	}
	cond := db.And(
		db.Cond{`customer_id`: flow.CustomerId},
		db.Cond{`asset_type`: flow.AssetType},
		db.Cond{`amount_type`: flow.AmountType},
		db.Cond{`source_type`: flow.SourceType},
		db.Cond{`source_table`: flow.SourceTable},
		db.Cond{`source_id`: flow.SourceId},
	)

	m := dbschema.NewOfficialCustomerWalletFlow(f.Context())
	maxN, err := m.Param(nil, cond.And()).Stat(`MAX`, `number`)
	if err != nil {
		return err
	}
	flow.Number = uint64(maxN) + 1
	return f.AddFlow(flow)
}

// AddFlow 添加不允许重复的流水记录
func (f *Wallet) AddFlow(flows ...*dbschema.OfficialCustomerWalletFlow) (err error) {
	flow := f.Flow
	if len(flows) > 0 {
		flow = flows[0]
	}
	err = f.Context().Begin()
	if err != nil {
		return
	}
	cond := db.And(
		db.Cond{`customer_id`: flow.CustomerId},
		db.Cond{`asset_type`: flow.AssetType},
	)
	var exists bool
	exists, err = f.Exists(nil, cond)
	if err != nil {
		f.Context().Rollback()
		return err
	}
	if !exists {
		f.Freeze = 0
		f.Balance = 0
		if flow.AmountType == `balance` { // 余额操作
			if flow.Amount < 0 { // 扣款操作
				f.Context().Rollback()
				return common.ErrBalanceNoEnough.SetZone(`balance`)
			}
			// 加款操作
			f.Balance = flow.Amount
			flow.WalletAmount = f.Balance
			if !AssetTypeIsIgnoreAccumulated(flow.AssetType) {
				f.Accumulated = flow.Amount
			}
		} else { // 冻结金额操作
			if flow.Amount < 0 { // 扣除冻结
				f.Context().Rollback()
				return f.Context().NewError(code.BalanceNoEnough, `冻结额不能小于0`).SetZone(`freeze`)
			}
			// 增加冻结
			f.Freeze = flow.Amount
			flow.WalletAmount = f.Freeze
		}
		f.CustomerId = flow.CustomerId
		f.AssetType = flow.AssetType
		_, err = f.Insert()
	} else {
		err = xdatabase.GetAndLock(f.OfficialCustomerWallet, cond)
		if err != nil {
			f.Context().Rollback()
			return err
		}
		incrAmount := param.AsString(flow.Amount)
		kvset := echo.H{
			flow.AmountType: db.Raw(flow.AmountType + `+` + incrAmount),
		}
		var amount float64
		if flow.AmountType == `balance` {
			amount = f.Balance
			if flow.Amount > 0 && !AssetTypeIsIgnoreAccumulated(flow.AssetType) { // 加款操作时增加累计加款金额
				kvset.Set(`accumulated`, db.Raw(`accumulated+`+incrAmount))
			}
		} else {
			amount = f.Freeze
		}
		//flow.WalletAmount = amount+flow.Amount
		oldAmountDecimal := decimal.NewFromFloat(amount)
		newAmountDecimal := decimal.NewFromFloat(flow.Amount)
		newAmountDecimal = oldAmountDecimal.Add(newAmountDecimal).Truncate(2)
		flow.WalletAmount, _ = newAmountDecimal.Float64()
		if flow.WalletAmount < 0 { //处理flow.Amount为负数(即扣除余额)的情况。扣款操作时检查余额是否足够
			f.Context().Rollback()
			if sessdata.User(f.Context()) != nil {
				return f.Context().NewError(code.BalanceNoEnough, `扣除余额(%v)失败！客户(ID:%d)的余额不足`, flow.Amount, flow.CustomerId).SetZone(`balance`)
			}
			return common.ErrBalanceNoEnough
		}
		err = f.UpdateFields(nil, kvset, cond)
	}
	if err != nil {
		f.Context().Rollback()
		return
	}
	/*
		flow := &dbschema.OfficialCustomerWalletFlow{
			CustomerId:  0,
			SourceCustomer: 0,
			AssetType:  `money`,
			AmountType: AmountTypeBalance, //金额类型(balance-余额;freeze-冻结额)
			Amount:      0,
			SourceType:  ``,
			SourceTable: ``,
			SourceId:    0,
			TradeNo:     ``,
			Status:      FlowStatusPending, //状态(pending-待确认;confirmed-已确认;canceled-已取消)
			Description: ``,
		}
	*/
	_, err = flow.Insert()
	f.Context().End(err == nil)
	return
}

func (f *Wallet) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*WalletExt, error) {
	list := []*WalletExt{}
	_, err := common.NewLister(f.OfficialCustomerWallet, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...).Relation(`Customer`, CusomterSafeFieldsSelector)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return list, err
	}
	for k, v := range list {
		v.AssetTypeName = AssetTypes.Get(v.AssetType)
		list[k] = v
	}
	return list, err
}

func (f *Wallet) FlowListPage(cond *db.Compounds, orderby ...interface{}) ([]*WalletFlowExt, error) {
	list := []*WalletFlowExt{}
	_, err := common.NewLister(f.Flow, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...).Relation(`Customer`, CusomterSafeFieldsSelector).Relation(`SourceCustomer`, CusomterSafeFieldsSelector)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return list, err
	}
	for k, v := range list {
		v.AssetTypeName = AssetTypes.Get(v.AssetType)
		list[k] = v
	}
	return list, err
}

var (
	emptyAsset  = dbschema.NewOfficialCustomerWallet(nil)
	emptyCredit = &dbschema.OfficialCustomerWallet{Balance: 10}
)

func (f *Wallet) GetAssetBalance(assetType string, customerID uint64) (float64, error) {
	recv := dbschema.NewOfficialCustomerWallet(f.Context())
	err := recv.Get(func(r db.Result) db.Result {
		return r.Select(`balance`)
	}, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`asset_type`: assetType},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
	}
	return recv.Balance, err
}

func (f *Wallet) GetAsset(assetType string, customerID uint64) (*dbschema.OfficialCustomerWallet, error) {
	recv := dbschema.NewOfficialCustomerWallet(f.Context())
	err := recv.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`asset_type`: assetType},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
	}
	return recv, err
}

// ListCustomerAllAssets 列出客户所有资产
func (f *Wallet) ListCustomerAllAssets(customerID uint64) (map[string]*dbschema.OfficialCustomerWallet, error) {
	var assetTypes []string
	assets := map[string]*dbschema.OfficialCustomerWallet{}
	assetTypeList := AssetTypeList()
	for _, assetType := range assetTypeList {
		if assetType.K == `credit` {
			assets[assetType.K] = emptyCredit
		} else {
			assets[assetType.K] = emptyAsset
		}
		assetTypes = append(assetTypes, assetType.K)
	}
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`asset_type`: db.In(assetTypes)},
	))
	if err != nil {
		return nil, err
	}
	for _, asset := range f.Objects() {
		assets[asset.AssetType] = asset
	}
	return assets, err
}
