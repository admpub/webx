package xrole

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

type (
	AuthChecker func(
		c echo.Context,
		rpath string,
		customer *dbschema.OfficialCustomer,
		permission *RolePermission,
	) (err error, ppath string, returning bool)
	AuthCheckers map[string]AuthChecker
)

func (a AuthCheckers) Check(
	c echo.Context,
	rpath string,
	customer *dbschema.OfficialCustomer,
	permission *RolePermission,
) (err error, ppath string, returning bool) {
	if checker, ok := a[rpath]; ok {
		return checker(c, rpath, customer, permission)
	}
	ppath = rpath
	return
}

var SpecialAuths = AuthCheckers{
	`/user/file/crop`: func(
		c echo.Context,
		rpath string,
		customer *dbschema.OfficialCustomer,
		permission *RolePermission,
	) (err error, ppath string, returning bool) {
		ppath = `/user/file/upload/:type`
		return
	},
}

func AuthRegister(ppath string, checker AuthChecker) {
	SpecialAuths[ppath] = checker
}

func AuthUnregister(ppath string) {
	if _, ok := SpecialAuths[ppath]; ok {
		delete(SpecialAuths, ppath)
	}
}
