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

package middleware

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/mwutils"
	"github.com/admpub/webx/application/middleware/sessdata"
	"github.com/admpub/webx/application/model/official"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

var TmplFuncGenerator = mwutils.TmplFuncGenerators{}

func SetFunc(ctx echo.Context) error {
	ctx.SetFunc(`Customer`, func() *dbschema.OfficialCustomer {
		return Customer(ctx)
	})
	ctx.SetFunc(`CustomerDetail`, func() *modelCustomer.CustomerAndGroup {
		return CustomerDetail(ctx)
	})
	ctx.SetFunc(`ImageProxyURL`, sessdata.ImageProxyURL)
	ctx.SetFunc(`ResizeImageURL`, sessdata.ResizeImageURL)
	ctx.SetFunc(`AbsoluteURL`, sessdata.AbsoluteURL)
	ctx.SetFunc(`PictureHTML`, sessdata.PictureWithDefaultHTML)
	ctx.SetFunc(`OutputContent`, sessdata.OutputContent)

	m := dbschema.NewOfficialCommonNavigate(ctx)

	ctx.SetFunc(`FrontendNav`, func(parentIDs ...uint) []*official.NavigateExt {
		return navigateList(ctx, m, `default`, parentIDs...)
	})

	ctx.SetFunc(`CustomerNav`, func(parentIDs ...uint) []*official.NavigateExt {
		return navigateList(ctx, m, `userCenter`, parentIDs...)
	})

	ctx.SetFunc(`Friendlink`, func(limit int, categoryIds ...uint) []*dbschema.OfficialCommonFriendlink {
		m := official.NewFriendlink(ctx)
		list, _ := m.ListShowAndRecord(limit, categoryIds...)
		return list
	})

	TmplFuncGenerator.Apply(ctx)
	return nil
}

func CustomerDetail(c echo.Context) *modelCustomer.CustomerAndGroup {
	customerDetail, _ := c.Internal().Get(`customerDetail`).(*modelCustomer.CustomerAndGroup)
	return customerDetail
}

func navigateList(ctx echo.Context, m *dbschema.OfficialCommonNavigate, navType string, parentIDs ...uint) []*official.NavigateExt {
	internalKey := `navigate.` + navType
	nav, ok := ctx.Internal().Get(internalKey).([]*official.NavigateExt)
	if !ok {
		nav = []*official.NavigateExt{}
		m.ListByOffset(&nav, func(r db.Result) db.Result {
			return r.OrderBy(`level`, `sort`, `id`)
		}, 0, -1, db.And(
			db.Cond{`disabled`: `N`},
			db.Cond{`type`: navType},
		))
		ctx.Internal().Set(internalKey, nav)
	}
	if len(parentIDs) > 0 {
		key := internalKey + `.` + fmt.Sprint(parentIDs[0])
		navList, ok := ctx.Internal().Get(key).([]*official.NavigateExt)
		if !ok {
			navList = []*official.NavigateExt{}
			for _, v := range nav {
				if v.ParentId == parentIDs[0] {
					navList = append(navList, v)
				}
			}
			ctx.Internal().Set(key, navList)
		}
		return navList
	}
	return nav
}

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			SetFunc(c)
			return h.Handle(c)
		})
	}
}
