package customer

import (
	"github.com/admpub/decimal"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/model/base"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/middleware/sessdata"
)

func NewWallet(ctx echo.Context) *Wallet {
	m := &Wallet{
		Wallet: dbschema.NewOfficialCustomerWallet(ctx),
		Flow:   dbschema.NewOfficialCustomerWalletFlow(ctx),
		base:   base.New(ctx),
	}
	return m
}

type Wallet struct {
	Wallet *dbschema.OfficialCustomerWallet
	Flow   *dbschema.OfficialCustomerWalletFlow
	base   *base.Base
}

func (f *Wallet) Use(tx *factory.Transaction) *Wallet {
	f.Wallet.Use(tx)
	f.Flow.Use(tx)
	return f
}

func (f *Wallet) Add() (pk interface{}, err error) {
	return f.Wallet.Insert()
}

func (f *Wallet) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.Wallet.Update(mw, args...)
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

	m := dbschema.NewOfficialCustomerWalletFlow(f.base.Context)
	m.Use(f.base.Tx())
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
	f.Wallet.Use(f.base.Tx())
	flow.Use(f.base.Tx())
	cond := db.And(
		db.Cond{`customer_id`: flow.CustomerId},
		db.Cond{`asset_type`: flow.AssetType},
	)
	if err := f.base.Context.Begin(); err != nil {
		return err
	}
	err = f.Wallet.Get(nil, cond)
	if err != nil {
		if err != db.ErrNoMoreRows {
			f.base.Context.Rollback()
			return err
		}
		f.Wallet.Freeze = 0
		f.Wallet.Balance = 0
		if flow.AmountType == `balance` { // 余额操作
			if flow.Amount < 0 { // 扣款操作
				f.base.Context.Rollback()
				return common.ErrBalanceNoEnough
			}
			// 加款操作
			f.Wallet.Balance = flow.Amount
			flow.WalletAmount = f.Wallet.Balance
		} else { // 冻结金额操作
			if flow.Amount < 0 { // 扣除冻结
				f.base.Context.Rollback()
				return f.base.E(`冻结额不能小于0`)
			}
			// 增加冻结
			f.Wallet.Freeze = flow.Amount
			flow.WalletAmount = f.Wallet.Freeze
		}
		f.Wallet.CustomerId = flow.CustomerId
		f.Wallet.AssetType = flow.AssetType
		_, err = f.Wallet.Insert()
	} else {
		var amount float64
		if flow.AmountType == `balance` {
			amount = f.Wallet.Balance
		} else {
			amount = f.Wallet.Freeze
		}
		//flow.WalletAmount = amount+flow.Amount
		oldAmountDecimal := decimal.NewFromFloat(amount)
		newAmountDecimal := decimal.NewFromFloat(flow.Amount)
		newAmountDecimal = oldAmountDecimal.Add(newAmountDecimal).Truncate(2)
		flow.WalletAmount, _ = newAmountDecimal.Float64()
		if flow.WalletAmount < 0 { //处理flow.Amount为负数(即扣除余额)的情况。扣款操作时检查余额是否足够
			f.base.Context.Rollback()
			if sessdata.User(f.base.Context) != nil {
				return f.base.E(`扣除余额(%v)失败！客户(ID:%d)的余额不足`, flow.Amount, flow.CustomerId)
			}
			return common.ErrBalanceNoEnough
		}
		incrAmount := param.AsString(flow.Amount)
		err = f.Wallet.UpdateField(nil, flow.AmountType, db.Raw(flow.AmountType+`+`+incrAmount), cond)
	}
	if err != nil {
		f.base.Context.Rollback()
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
	f.base.Context.End(err == nil)
	return
}

func (f *Wallet) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*WalletExt, error) {
	list := []*WalletExt{}
	_, err := common.NewLister(f.Wallet, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...).Relation(`Customer`, CusomterSafeFieldsSelector)
	}, cond.And()).Paging(f.base.Context)
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
	}, cond.And()).Paging(f.base.Context)
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
	recv := dbschema.NewOfficialCustomerWallet(f.base.Context)
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
	recv := dbschema.NewOfficialCustomerWallet(f.base.Context)
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
	_, err := f.Wallet.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`asset_type`: db.In(assetTypes)},
	))
	if err != nil {
		return nil, err
	}
	for _, asset := range f.Wallet.Objects() {
		assets[asset.AssetType] = asset
	}
	return assets, err
}
