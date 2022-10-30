/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package role

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/library/xrole"
	"github.com/admpub/webx/application/library/xrole/xroleutils"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func Index(ctx echo.Context) error {
	m := modelCustomer.NewRole(ctx)
	_, err := handler.PagingWithLister(ctx, m)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`official/customer/role/index`, handler.Err(ctx, err))
}

func Add(ctx echo.Context) error {
	var err error
	m := modelCustomer.NewRole(ctx)
	permission := xrole.NewRolePermission()
	if ctx.IsPost() {
		ctx.Begin()
		err = ctx.MustBind(m.OfficialCustomerRole)
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			err = xroleutils.AddCustomerRolePermission(ctx, m.Id)
		}
		ctx.End(err == nil)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/official/customer/role/index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCustomerRole, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
				rpM := modelCustomer.NewRolePermission(ctx)
				rpM.ListByOffset(nil, nil, 0, -1, `role_id`, m.Id)
				permissionList := []*xrole.CustomerRoleWithPermissions{
					{
						OfficialCustomerRole: m.OfficialCustomerRole,
						Permissions:          rpM.Objects(),
					},
				}
				permission.Init(permissionList)
			}
		}
	}
	ctx.Set(`activeURL`, `/official/customer/role/index`)
	ctx.Set(`data`, m)
	ctx.Set(`permission`, permission)
	ctx.Set(`permissionTypes`, xrole.CustomerRolePermissionType.Slice())
	xroleutils.CustomerRolePermissionTypeFireRender(ctx)
	return ctx.Render(`official/customer/role/edit`, handler.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewRole(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/official/customer/role/index`))
	}
	if ctx.IsPost() {
		ctx.Begin()
		err = ctx.MustBind(m.OfficialCustomerRole)
		if err == nil {
			m.Id = id
		}
		if err == nil {
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			err = xroleutils.EditCustomerRolePermission(ctx, m.Id)
		}
		ctx.End(err == nil)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/official/customer/role/index`))
		}
	}

	echo.StructToForm(ctx, m.OfficialCustomerRole, ``, echo.LowerCaseFirstLetter)
	ctx.Set(`activeURL`, `/official/customer/role/index`)
	ctx.Set(`data`, m)
	rpM := modelCustomer.NewRolePermission(ctx)
	rpM.ListByOffset(nil, nil, 0, -1, `role_id`, m.Id)
	permissionList := []*xrole.CustomerRoleWithPermissions{
		{
			OfficialCustomerRole: m.OfficialCustomerRole,
			Permissions:          rpM.Objects(),
		},
	}
	permission := xrole.NewRolePermission().Init(permissionList)
	ctx.Set(`permission`, permission)
	ctx.Set(`permissionTypes`, xrole.CustomerRolePermissionType.Slice())
	xroleutils.CustomerRolePermissionTypeFireRender(ctx)
	return ctx.Render(`official/customer/role/edit`, handler.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewRole(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		rpM := modelCustomer.NewRolePermission(ctx)
		rpM.Delete(nil, `role_id`, id)
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/customer/role/index`))
}
