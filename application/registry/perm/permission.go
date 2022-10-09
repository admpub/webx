package perm

import (
	"github.com/admpub/webx/application/library/xrole"
)

func New() *RolePermission {
	return xrole.NewRolePermission()
}

type RolePermission = xrole.RolePermission
