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

package customer

import (
	"strings"

	"github.com/admpub/nging/v4/application/library/perm"
	"github.com/admpub/webx/application/dbschema"
	modelLevel "github.com/admpub/webx/application/model/official/level"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func NewRole(ctx echo.Context) *Role {
	m := &Role{
		OfficialCustomerRole: dbschema.NewOfficialCustomerRole(ctx),
	}
	return m
}

type Role struct {
	*dbschema.OfficialCustomerRole
	permActions   *perm.Map
	permBehaviors perm.BehaviorPerms
}

func (u *Role) check() error {
	if len(u.Name) == 0 {
		return u.Context().NewError(code.InvalidParameter, `角色名不能为空`)
	}
	var exists bool
	var err error
	if u.Id > 0 {
		exists, err = u.Exists2(u.Name, u.Id)
	} else {
		exists, err = u.Exists(u.Name)
	}
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `角色名已经存在`)
	}
	return err
}

func (u *Role) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	if u.IsDefault == `Y` {
		if err := u.CancelDefault(); err != nil {
			return nil, err
		}
	}
	return u.OfficialCustomerRole.Insert()
}

func (u *Role) CancelDefault(excludeID ...uint) error {
	cond := db.NewCompounds()
	cond.Add(db.Cond{`is_default`: `Y`})
	if len(excludeID) > 0 && excludeID[0] > 0 {
		cond.Add(db.Cond{`id`: db.NotEq(excludeID[0])})
	}
	return u.UpdateFields(nil, echo.H{`is_default`: `N`}, cond.And())
}

func (u *Role) GetDefault() error {
	cond := db.NewCompounds()
	cond.Add(db.Cond{`disabled`: `N`})
	cond.Add(db.Cond{`is_default`: `Y`})
	return u.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And())
}

func (u *Role) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := u.check(); err != nil {
		return err
	}
	if u.IsDefault == `Y` {
		if err := u.CancelDefault(u.Id); err != nil {
			return err
		}
	}
	return u.OfficialCustomerRole.Update(mw, args...)
}

func (u *Role) Exists(name string) (bool, error) {
	return u.OfficialCustomerRole.Exists(nil, db.Cond{`name`: name})
}

func (u *Role) ListRoleIDsByCustomer(customer *dbschema.OfficialCustomer) (roleIDs []uint) {
	// 客户自身的角色
	rolesString := customer.RoleIds

	// 客户所属等级的角色
	levelM := modelLevel.NewLevel(u.Context())
	levelList, err := levelM.ListByCustomer(customer)
	if err == nil {
		for _, level := range levelList {
			if len(level.RoleIds) > 0 {
				rolesString += `,` + level.RoleIds
			}
		}
	}
	if customer.LevelId > 0 {
		err = levelM.Get(nil, db.And(
			db.Cond{`id`: customer.LevelId},
			db.Cond{`disabled`: `N`},
		))
		if err == nil {
			if len(levelM.RoleIds) > 0 {
				rolesString += `,` + levelM.RoleIds
			}
		}
	}

	// 客户代理等级的角色
	if customer.AgentLevel > 0 {
		levelM := dbschema.NewOfficialCustomerAgentLevel(u.Context())
		err = levelM.Get(nil, db.And(
			db.Cond{`id`: customer.AgentLevel},
			db.Cond{`disabled`: `N`},
		))
		if err == nil {
			if len(levelM.RoleIds) > 0 {
				rolesString += `,` + levelM.RoleIds
			}
		}
	}

	// 提取角色ID(自动跳过重复ID)
	rolesString = strings.Trim(rolesString, `,`)
	if len(rolesString) > 0 {
		temp := map[uint]struct{}{}
		roleIDs = param.StringSlice(strings.Split(rolesString, `,`)).Uint(func(_ int, v uint) bool {
			if v == 0 {
				return false
			}
			if _, ok := temp[v]; ok {
				return false
			}
			temp[v] = struct{}{}
			return true
		})
		temp = nil
	}
	return
}

func (u *Role) ListByCustomer(customer *dbschema.OfficialCustomer) (roleList []*dbschema.OfficialCustomerRole) {
	roleIDs := u.ListRoleIDsByCustomer(customer)
	// 获取角色数据
	if len(roleIDs) > 0 {
		u.ListByOffset(nil, nil, 0, -1, db.And(
			db.Cond{`disabled`: `N`},
			db.Cond{`id`: db.In(roleIDs)},
		))
		roleList = u.Objects()
	}
	return
}

func (u *Role) Exists2(name string, excludeID uint) (bool, error) {
	n, e := u.Param(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(excludeID)},
	)).Count()
	return n > 0, e
}
