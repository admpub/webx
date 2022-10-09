package xroleutils

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/perm"
	"github.com/admpub/webx/application/library/xrole"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func CustomerRolePermissionTypeFireRender(ctx echo.Context) (err error) {
	return perm.HandleFireRender(ctx, xrole.CustomerRolePermissionType)
}

func CustomerRolePermissionTypeGenerate(ctx echo.Context) (mp map[string]string, err error) {
	return perm.HandleGenerate(ctx, xrole.CustomerRolePermissionType)
}

func CustomerRolePermissionTypeCheck(ctx echo.Context, current string, typ string, permission string, parsed interface{}) (mp interface{}, err error) {
	return perm.HandleCheck(ctx, xrole.CustomerRolePermissionType, current, typ, permission, parsed)
}

func AddCustomerRolePermission(ctx echo.Context, roleID uint) (err error) {
	var perms map[string]string
	perms, err = CustomerRolePermissionTypeGenerate(ctx)
	if err != nil {
		return
	}
	rpM := modelCustomer.NewRolePermission(ctx)
	for typ, perm := range perms {
		rpM.RoleId = roleID
		rpM.Type = typ
		rpM.Permission = perm
		_, err = rpM.Add()
		if err != nil {
			break
		}
	}
	return
}

func EditCustomerRolePermission(ctx echo.Context, roleID uint) (err error) {
	var perms map[string]string
	perms, err = CustomerRolePermissionTypeGenerate(ctx)
	if err != nil {
		return
	}
	rpM := modelCustomer.NewRolePermission(ctx)
	_, err = rpM.ListByOffset(nil, nil, 0, -1, db.Cond{`role_id`: roleID})
	if err != nil {
		return
	}
	var deleted []string
	for _, rule := range rpM.Objects() {
		_, ok := perms[rule.Type]
		if !ok {
			deleted = append(deleted, rule.Type)
		}
	}
	if len(deleted) > 0 {
		err = rpM.Delete(nil, db.And(
			db.Cond{`role_id`: roleID},
			db.Cond{`type`: db.In(deleted)},
		))
		if err != nil {
			return
		}
	}
	for typ, perm := range perms {
		rpM.RoleId = roleID
		rpM.Type = typ
		rpM.Permission = perm
		_, err = rpM.Add()
		if err != nil {
			break
		}
	}
	return
}
