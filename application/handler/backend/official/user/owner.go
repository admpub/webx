package user

import (
	"github.com/admpub/null"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webcore/model"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

func SelectCustomer(ctx echo.Context) error {
	m := modelCustomer.NewCustomer(ctx)
	cond := &db.Compounds{}
	nsql.SelectPageCond(ctx, cond, `id`, `name%`)
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`).OrderBy(`-id`)
	}, cond.And()))
	if err != nil {
		return err
	}
	ctx.Set(`listData`, m.Objects())
	return ctx.JSON(ctx.Data().SetData(ctx.Stored()))
}

func SelectUser(ctx echo.Context) error {
	cond := &db.Compounds{}
	nsql.SelectPageCond(ctx, cond, `id`, `username%`)
	m := model.NewUser(ctx)
	listData := []null.StringMap{}
	_, err := common.PagingWithLister(ctx, common.NewLister(m, &listData, func(r db.Result) db.Result {
		return r.Select(`id`, db.Raw(`username AS name`)).OrderBy(`-id`)
	}, cond.And()))
	if err != nil {
		return err
	}
	ctx.Set(`listData`, listData)
	return ctx.JSON(ctx.Data().SetData(ctx.Stored()))
}

func SelectOwner(ctx echo.Context, operation string) error {
	if len(operation) > 0 {
		switch operation {
		case `customer`:
			return SelectCustomer(ctx)
		case `user`:
			return SelectUser(ctx)
		}
	}
	return ctx.JSON(ctx.Data().SetData(echo.H{`listData`: nil, `pagination`: pagination.New(ctx)}))
}
