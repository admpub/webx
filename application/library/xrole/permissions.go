package xrole

import (
	"github.com/admpub/nging/v4/application/library/role"
	"github.com/admpub/webx/application/dbschema"
)

type CustomerRoleWithPermissions struct {
	*dbschema.OfficialCustomerRole
	Permissions []*dbschema.OfficialCustomerRolePermission `db:"-,relation=role_id:id"`
}

func (u *CustomerRoleWithPermissions) GetPermissions() []role.PermissionConfiger {
	r := make([]role.PermissionConfiger, len(u.Permissions))
	for k, v := range u.Permissions {
		r[k] = v
	}
	return r
}
