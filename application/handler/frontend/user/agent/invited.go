package agent

import (
	"github.com/admpub/nging/v5/application/handler"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

// InvitedList 被邀请的用户列表
func InvitedList(ctx echo.Context) error {
	m := modelCustomer.NewCustomer(ctx)
	customer := xMW.Customer(ctx)
	cond := db.Cond{
		`inviter_id`: customer.Id,
	}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m.OfficialCustomer, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond))
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`activeURL`, `/user/agent/index`)
	ctx.Set(`title`, ctx.T(`我邀请的用户`))
	return ctx.Render(`/user/agent/invited`, handler.Err(ctx, err))
}
