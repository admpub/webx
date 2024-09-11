package membership

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	modelLevel "github.com/admpub/webx/application/model/official/level"
	"github.com/coscms/webcore/library/common"
)

func Index(ctx echo.Context) error {
	m := modelCustomer.NewGroupPackage(ctx)
	groups, err := m.ListGroup()
	if err != nil {
		return err
	}
	group := ctx.Form(`group`)
	if len(group) > 0 {
		if !com.InSlice(group, groups) {
			return ctx.NewError(code.InvalidParameter, `无效的组: %s`, group).SetZone(`group`)
		}
	} else if len(groups) > 0 {
		group = groups[0]
	}
	customer := sessdata.Customer(ctx)
	var packageList []*dbschema.OfficialCustomerGroupPackage
	var myLevel *modelLevel.RelationExt
	if len(group) > 0 {
		err = m.ListByGroup(group)
		if err != nil {
			return err
		}
		packageList = m.Objects()
		levelM := modelCustomer.NewLevel(ctx)
		myLevel, err = levelM.GetByCustomerID(group, customer.Id)
		if err != nil {
			if err != db.ErrNoMoreRows {
				return err
			}
			err = nil
		}
	}
	groupList := make([]echo.KV, len(groups))
	for i, v := range groups {
		item := modelLevel.GroupList.GetItem(v)
		if item != nil {
			groupList[i] = *item
		}
	}
	ctx.Set(`groupList`, groupList)
	ctx.Set(`packageList`, packageList)
	ctx.Set(`group`, group)
	ctx.Set(`myLevel`, myLevel)
	ctx.SetFunc(`timeUnitSuffix`, func(n uint, unit string) string {
		return modelCustomer.GroupPackageTimeUnitSuffix(ctx, n, unit)
	})
	return ctx.Render(`user/membership/index`, common.Err(ctx, err))
}

func Buy(ctx echo.Context) error {
	packageID := ctx.Paramx(`packageId`).Uint()
	if packageID < 1 {
		return ctx.NewError(code.InvalidParameter, `参数无效`).SetZone(`packageId`)
	}
	pkgM := modelCustomer.NewGroupPackage(ctx)
	err := pkgM.Get(nil, `id`, packageID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return ctx.NewError(code.DataNotFound, `套餐不存在`).SetZone(`packageId`)
		}
		return err
	}
	customer := sessdata.Customer(ctx)
	myLevelM := modelCustomer.NewLevel(ctx)
	var myLevel *modelLevel.RelationExt
	myLevel, err = myLevelM.GetByCustomerID(pkgM.Group, customer.Id)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		err = nil
	} else {
		if myLevel.Expired == 0 {
			return ctx.NewError(code.Failure, `您已经是终身“%s”，无需再次购买`, modelLevel.GroupList.Get(pkgM.Group))
		}
	}
	m := modelCustomer.NewWallet(ctx)
	var money float64
	money, err = m.GetBalance(`money`, customer.Id)
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		now := time.Now()
		var baseTime time.Time
		if myLevel != nil {
			if myLevel.Expired == 0 {
				return ctx.NewError(code.Failure, `您已经是终身“%s”，无需再次购买`, modelLevel.GroupList.Get(pkgM.Group))
			}
			baseTime = time.Unix(int64(myLevel.Expired), 0)
			if baseTime.Before(now) {
				baseTime = now
			}
		}
		expiresTime := pkgM.MakeExpireTime(baseTime)
		myLevelM.CustomerId = customer.Id
		if myLevel != nil {
			myLevelM.LevelId = myLevel.LevelId
		} else {
			lvM := modelLevel.NewLevel(ctx)
			err = lvM.GetMinLevelByGroup(pkgM.Group)
			if err != nil {
				if err == db.ErrNoMoreRows {
					return ctx.NewError(code.DataNotFound, `“%s”尚未配置等级，暂时无法购买`, modelLevel.GroupList.Get(pkgM.Group))
				}
				return err
			}
			myLevelM.LevelId = lvM.Id
		}
		myLevelM.Expired = uint(expiresTime.Unix())
		myLevelM.AccumulatedDays = (myLevelM.Expired - uint(baseTime.Unix())) / 86400
		ctx.Begin()
		// 钱包余额支付
		walletM := modelCustomer.NewWallet(ctx)
		walletM.Flow.CustomerId = customer.Id
		walletM.Flow.AssetType = modelCustomer.AssetTypeMoney
		walletM.Flow.AmountType = modelCustomer.AmountTypeBalance
		walletM.Flow.Amount = -pkgM.Price
		walletM.Flow.SourceType = `buy`
		walletM.Flow.SourceTable = `official_customer_group_package`
		walletM.Flow.SourceId = uint64(pkgM.Id)
		walletM.Flow.TradeNo = ``
		walletM.Flow.Status = modelCustomer.FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
		walletM.Flow.Description = `购买会员套餐: ` + pkgM.Title
		err = walletM.AddFlow()
		if err != nil {
			ctx.Rollback()
			goto END
		}

		// 添加会员等级数据
		_, err = myLevelM.Add()
		if err != nil {
			ctx.Rollback()
			goto END
		}
		err = pkgM.IncrSold(pkgM.Id)
		if err != nil {
			ctx.Rollback()
			goto END
		}
		ctx.Commit()
		next := ctx.Form(`next`)
		if len(next) == 0 {
			next = sessdata.URLFor(`/user/membership/index?group=` + pkgM.Group)
		}
		ctx.Data().SetInfo(ctx.T(`购买成功`))
		return ctx.Redirect(next)
	}

END:
	ctx.Set(`package`, pkgM.OfficialCustomerGroupPackage)
	ctx.Set(`myLevel`, myLevel)
	ctx.Set(`money`, money)
	ctx.Set(`activeURL`, `/user/membership/index`)
	ctx.SetFunc(`timeUnitSuffix`, func(n uint, unit string) string {
		return modelCustomer.GroupPackageTimeUnitSuffix(ctx, n, unit)
	})
	return ctx.Render(`user/membership/buy`, common.Err(ctx, err))
}
