package membership

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/dbschema"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func Index(ctx echo.Context) error {
	m := modelCustomer.NewGroupPackage(ctx)
	groupList, err := m.ListGroup()
	if err != nil {
		return err
	}
	ctx.Set(`groupList`, groupList)
	var packageList []*dbschema.OfficialCustomerGroupPackage
	if len(groupList) > 0 {
		err = m.ListByGroup(groupList[0])
		if err != nil {
			return err
		}
		packageList = m.Objects()
	}
	ctx.Set(`packageList`, packageList)
	ctx.SetFunc(`timeUnitName`, modelCustomer.GroupPackageTimeUnits.Get)
	return ctx.Render(`user/membership/index`, handler.Err(ctx, err))
}

func Buy(ctx echo.Context) error {
	return ctx.Render(`user/membership/buy`, nil)
}
