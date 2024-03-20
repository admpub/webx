package customer

import (
	"fmt"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/cache"
	"github.com/admpub/webx/application/library/xrole"
	"github.com/admpub/webx/application/middleware/sessdata"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func CustomerPermTTL(c echo.Context) int64 {
	cacheTTL, ok := c.Internal().Get(`customerPermCacheTTL`).(int64)
	if ok {
		return cacheTTL
	}
	cacheTTL = config.Setting(`base`).Int64(`customerPermCacheTTL`, cache.CacheDisabled)
	if cacheTTL == 0 {
		cacheTTL = int64(cache.CacheDisabled)
	} else if cacheTTL < int64(cache.CacheDisabled) {
		cacheTTL = int64(cache.CacheFresh)
	}
	c.Internal().Set(`customerPermCacheTTL`, cacheTTL)
	return cacheTTL
}

func CustomerPermission(c echo.Context, customers ...*dbschema.OfficialCustomer) *xrole.RolePermission {
	permission, ok := c.Internal().Get(`customerPermission`).(*xrole.RolePermission)
	if !ok || permission == nil {
		var customer *dbschema.OfficialCustomer
		if len(customers) > 0 && customers[0] != nil {
			customer = customers[0]
		} else {
			customer = sessdata.Customer(c)
		}
		if customer == nil {
			return nil
		}
		customerID := fmt.Sprint(customer.Id)
		permission = xrole.NewRolePermission()
		cache.XFunc(c, sessdata.PermissionCacheKey+customerID, permission, func() error {
			permission.Init(CustomerRoles(c, customer))
			return nil
		}, cache.GetTTLByNumber(CustomerPermTTL(c), nil))
		c.Internal().Set(`customerPermission`, permission)
	}
	return permission
}

func CustomerRoles(c echo.Context, customers ...*dbschema.OfficialCustomer) (roleList []*xrole.CustomerRoleWithPermissions) {
	roleList, ok := c.Internal().Get(`customerRoles`).([]*xrole.CustomerRoleWithPermissions)
	if ok {
		return roleList
	}
	var customer *dbschema.OfficialCustomer
	if len(customers) > 0 && customers[0] != nil {
		customer = customers[0]
	} else {
		customer = sessdata.Customer(c)
	}
	if customer == nil {
		return nil
	}
	roleM := NewRole(c)
	roleIDs := roleM.ListRoleIDsByCustomer(customer)
	if len(roleIDs) > 0 {
		roleM.ListByOffset(&roleList, nil, 0, -1, db.And(
			db.Cond{`disabled`: `N`},
			db.Cond{`id`: db.In(roleIDs)},
		))
	}
	c.Internal().Set(`customerRoles`, roleList)
	return roleList
}

func CustomerRolePermissionForBehavior(c echo.Context, behaviorName string, customer ...*dbschema.OfficialCustomer) interface{} {
	permission := CustomerPermission(c, customer...)
	//echo.Dump(permission)
	if permission == nil {
		return nil
	}
	bev, ok := permission.Get(c, xrole.CustomerRolePermissionTypeBehavior).(perm.BehaviorPerms)
	if !ok {
		return nil
	}
	//echo.Dump(bev)
	return bev.Get(behaviorName).Value
}
