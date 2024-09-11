package xrole

import (
	"github.com/coscms/webcore/library/role"
)

func NewRolePermission() *RolePermission {
	r := &RolePermission{}
	r.CommonPermission = role.NewCommonPermission(CustomerRolePermissionType, r)
	return r
}

type RolePermission struct {
	*role.CommonPermission
	Roles []*CustomerRoleWithPermissions
}

func (r *RolePermission) Init(roleList []*CustomerRoleWithPermissions) *RolePermission {
	r.Roles = roleList
	gts := make([]role.PermissionsGetter, len(r.Roles))
	for k, v := range r.Roles {
		gts[k] = v
	}
	r.CommonPermission.Init(gts)
	return r
}
