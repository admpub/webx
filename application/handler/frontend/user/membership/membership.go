package membership

import (
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	modelLevel "github.com/admpub/webx/application/model/official/level"
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
	return ctx.Render(`user/membership/index`, handler.Err(ctx, err))
}

func Buy(ctx echo.Context) error {
	return ctx.Render(`user/membership/buy`, nil)
}
