package xrole

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/perm"
	"github.com/coscms/webcore/library/role"
)

var CustomerRolePermissionType = echo.NewKVData().
	Add(CustomerRolePermissionTypePage, `页面权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`official/customer/role/edit_perm_page`).SetTmpl(`official/customer/role/edit_perm_page_foot`, `foot`).
			SetGenerator(PermPageGenerator).
			SetParser(PermPageParser).
			SetChecker(PermPageChecker).
			SetItemLister(PermPageList).
			OnRender(PermPageOnRender),
	)).
	Add(CustomerRolePermissionTypeBehavior, `行为权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`official/customer/role/edit_perm_behavior`).SetTmpl(`official/customer/role/edit_perm_behavior_foot`, `foot`).
			SetGenerator(PermBehaviorGenerator).
			SetParser(PermBehaviorParser).
			SetChecker(PermBehaviorChecker).
			SetItemLister(PermBehaviorList).
			OnRender(PermBehaviorOnRender).
			SetIsValid(PermBehaviorIsValid),
	))

const (
	CustomerRolePermissionTypePage     = role.RolePermissionTypePage
	CustomerRolePermissionTypeBehavior = role.RolePermissionTypeBehavior
)
