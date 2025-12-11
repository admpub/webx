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

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/library/xrole"
	"github.com/coscms/webfront/library/xrole/xroleutils"
	"github.com/coscms/webfront/model/i18nm"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func Index(ctx echo.Context) error {
	m := modelCustomer.NewRole(ctx)
	_, err := common.PagingWithLister(ctx, m)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`official/customer/role/index`, common.Err(ctx, err))
}

func Add(ctx echo.Context) error {
	var err error
	m := modelCustomer.NewRole(ctx)
	permission := xrole.NewRolePermission()
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialCustomerRole, uint64(id))
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
	form := formbuilder.New(ctx,
		m.OfficialCustomerRole,
		formbuilder.ConfigFile(`official/customer/role/edit`),
		formbuilder.AllowedNames(
			`parentId`, `disabled`, `isDefault`, `name`, `description`,
		),
	)
	form.OnPost(func() error {
		ctx.Begin()
		_, err := m.Add()
		if err != nil {
			ctx.Rollback()
			return err
		}
		err = xroleutils.AddCustomerRolePermission(ctx, m.Id)
		if err != nil {
			ctx.Rollback()
			return err
		}
		ctx.Commit()
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCustomerRole, uint64(m.Id), i18nm.OptionContentType(`description`, `text`))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/customer/role/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/official/customer/role/index`)
	ctx.Set(`data`, m)
	ctx.Set(`permission`, permission)
	ctx.Set(`permissionTypes`, xrole.CustomerRolePermissionType.Slice())
	xroleutils.CustomerRolePermissionTypeFireRender(ctx)
	return ctx.Render(`official/customer/role/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewRole(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		common.SendFail(ctx, err.Error())
		return ctx.Redirect(backend.URLFor(`/official/customer/role/index`))
	}
	if ctx.IsGet() {
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialCustomerRole, uint64(id))
	}
	form := formbuilder.New(ctx,
		m.OfficialCustomerRole,
		formbuilder.ConfigFile(`official/customer/role/edit`),
		formbuilder.AllowedNames(
			`parentId`, `disabled`, `isDefault`, `name`, `description`,
		),
	)
	form.OnPost(func() error {
		ctx.Begin()
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			ctx.Rollback()
			return err
		}
		err = xroleutils.EditCustomerRolePermission(ctx, m.Id)
		if err != nil {
			ctx.Rollback()
			return err
		}
		ctx.Commit()
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCustomerRole, uint64(m.Id), i18nm.OptionContentType(`description`, `text`))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(backend.URLFor(`/official/customer/role/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

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
	return ctx.Render(`official/customer/role/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewRole(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		rpM := modelCustomer.NewRolePermission(ctx)
		rpM.Delete(nil, `role_id`, id)
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/role/index`))
}
