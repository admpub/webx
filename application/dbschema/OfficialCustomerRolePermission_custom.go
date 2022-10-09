package dbschema

func (a *OfficialCustomerRolePermission) GetType() string {
	return a.Type
}

func (a *OfficialCustomerRolePermission) GetPermission() string {
	return a.Permission
}
