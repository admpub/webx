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

package invitation

import (
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func formFilter() echo.FormDataFilter {
	return formfilter.Build(
		formfilter.StartDateToTimestamp(`Start`),
		formfilter.EndDateToTimestamp(`End`),
		formfilter.JoinValues(`RoleIds`),
	)
}

func Index(ctx echo.Context) error {
	m := modelCustomer.NewInvitation(ctx)
	cond := db.Compounds{}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`code`, q)
	}
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := common.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`official/customer/invitation/index`, ret)
}

func Add(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := modelCustomer.NewInvitation(ctx)
		err = ctx.MustBind(m.OfficialCustomerInvitation, formFilter())
		if err == nil {
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			_, err = m.Add()
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/customer/invitation/index`))
		}
	} else {
		ctx.Request().Form().Set(`code`, com.RandomAlphanumeric(16))
	}
	ctx.Set(`activeURL`, `/official/customer/invitation/index`)
	roleM := modelCustomer.NewRole(ctx)
	roleM.List(nil, nil, 1, -1, `disabled`, `N`)
	ctx.Set(`roleList`, roleM.Objects())
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		return false
	})
	ctx.Set(`title`, ctx.T(`添加邀请码`))
	return ctx.Render(`official/customer/invitation/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewInvitation(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		common.SendFail(ctx, err.Error())
		return ctx.Redirect(backend.URLFor(`/official/customer/invitation/index`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerInvitation, formFilter())
		if err == nil {
			m.Id = id
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(backend.URLFor(`/official/customer/invitation/index`))
		}
	} else {
		echo.StructToForm(ctx, m.OfficialCustomerInvitation, ``, echo.LowerCaseFirstLetter)
		var startDate, endDate string
		if m.OfficialCustomerInvitation.Start > 0 {
			startDate = time.Unix(int64(m.OfficialCustomerInvitation.Start), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`start`, startDate)
		if m.OfficialCustomerInvitation.End > 0 {
			endDate = time.Unix(int64(m.OfficialCustomerInvitation.End), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`end`, endDate)
	}

	ctx.Set(`activeURL`, `/official/customer/invitation/index`)
	roleM := modelCustomer.NewRole(ctx)
	roleM.List(nil, nil, 1, -1, `disabled`, `N`)
	ctx.Set(`roleList`, roleM.Objects())
	var roleIds []uint
	if len(m.RoleIds) > 0 {
		roleIds = param.StringSlice(strings.Split(m.RoleIds, `,`)).Uint()
	}
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		for _, rid := range roleIds {
			if rid == roleId {
				return true
			}
		}
		return false
	})
	ctx.Set(`title`, ctx.T(`修改邀请码`))
	return ctx.Render(`official/customer/invitation/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewInvitation(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/invitation/index`))
}
