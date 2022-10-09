package bonus

import (
	"fmt"

	"github.com/admpub/webx/application/dbschema"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func New(ctx echo.Context) *Bonus {
	levelMdl := dbschema.NewOfficialCustomerAgentLevel(ctx)
	_, err := levelMdl.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: `N`})
	levelInfos := map[uint]*dbschema.OfficialCustomerAgentLevel{}
	if err == nil {
		for _, row := range levelMdl.Objects() {
			levelInfos[row.Id] = row
		}
	}
	return &Bonus{
		ctx:        ctx,
		parents:    map[uint64]struct{}{},
		levelInfos: levelInfos,
	}
}

const (
	MinAmount = 0.001
)

type Bonus struct {
	ctx                echo.Context
	levelInfos         map[uint]*dbschema.OfficialCustomerAgentLevel
	parents            map[uint64]struct{}
	currentLevelNumber int
	maxLevelNumber     int
	debug              bool
}

func (b *Bonus) SalesCommissionRatio(l *dbschema.OfficialCustomerAgentLevel, i int) float64 {
	switch i {
	case 1:
		return l.SalesCommissionRatio1
	case 2:
		return l.SalesCommissionRatio2
	case 3:
		return l.SalesCommissionRatio3
	case 4:
		return l.SalesCommissionRatio4
	case 5:
		return l.SalesCommissionRatio5
	case 6:
		return l.SalesCommissionRatio6
	case 7:
		return l.SalesCommissionRatio7
	case 8:
		return l.SalesCommissionRatio8
	case 9:
		return l.SalesCommissionRatio9
	case 10:
		return l.SalesCommissionRatio10
	default:
		return 0
	}
}

func (b *Bonus) SetDebug(on bool) *Bonus {
	b.debug = on
	return b
}

func (b *Bonus) log(format string, args ...interface{}) {
	if !b.debug {
		return
	}
	if len(args) > 0 {
		fmt.Println(format)
		return
	}
	fmt.Printf(format+"\n", args)
}

func (b *Bonus) Reset() *Bonus {
	b.currentLevelNumber = 0
	b.parents = map[uint64]struct{}{}
	b.levelInfos = map[uint]*dbschema.OfficialCustomerAgentLevel{}
	return b
}

func (b *Bonus) Confirm(customer *dbschema.OfficialCustomer, flow dbschema.OfficialCustomerWalletFlow) error {
	if customer.InviterId <= 0 { //没有邀请人
		b.log(`主动注册用户: %s(#%d)`, customer.Name, customer.Id)
		return nil
	}
	if _, ok := b.parents[customer.InviterId]; ok { //邀请链异常：出现重复邀请人
		b.log(ErrDeadLoop.Error()+`: %v`, b.parents)
		return ErrDeadLoop
	}
	b.parents[customer.InviterId] = struct{}{}
	b.currentLevelNumber++
	if b.currentLevelNumber > b.maxLevelNumber { //超过最大奖励层数
		return nil
	}
	parent := dbschema.NewOfficialCustomer(b.ctx)
	err := parent.Get(nil, `id`, customer.InviterId)
	if err != nil {
		return err
	}
	if parent.AgentLevel == 0 {
		b.log(`跳过非代理邀请人: %s(#%d)`, parent.Name, parent.Id)
		return err
	}
	levelInfo, ok := b.levelInfos[parent.AgentLevel]
	if !ok {
		b.log(ErrAgentLevelNotExists.Error()+`: %d`, parent.AgentLevel)
		return ErrAgentLevelNotExists
	}
	salesCommissionRatio := b.SalesCommissionRatio(levelInfo, b.currentLevelNumber)
	amount := salesCommissionRatio * flow.Amount
	if amount > MinAmount {
		walletM := modelCustomer.NewWallet(b.ctx)
		flowCopy := &dbschema.OfficialCustomerWalletFlow{
			CustomerId:     parent.Id,
			AssetType:      flow.AssetType, //资产类型(money-钱;point-点数;credit-信用分;integral-积分;gold-金币;silver-银币;copper-铜币)
			Amount:         amount,
			SourceCustomer: flow.CustomerId,
			SourceType:     flow.SourceType,
			SourceTable:    flow.SourceTable,
			SourceId:       flow.SourceId,
			TradeNo:        flow.TradeNo,
			Status:         flow.Status, //状态(pending-待确认;confirmed-已确认;canceled-已取消)
			Description:    flow.Description,
		}
		err = walletM.AddFlow(flowCopy)
		if err != nil {
			return err
		}
	}
	err = b.Confirm(parent, flow)
	return err
}
