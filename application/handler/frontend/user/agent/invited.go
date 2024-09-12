package agent

import (
	"github.com/coscms/webcore/library/common"
	xMW "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
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
	_, err := common.PagingWithLister(ctx, common.NewLister(m.OfficialCustomer, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond))
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`activeURL`, `/user/agent/index`)
	ctx.Set(`title`, ctx.T(`我邀请的用户`))
	return ctx.Render(`/user/agent/invited`, common.Err(ctx, err))
}
