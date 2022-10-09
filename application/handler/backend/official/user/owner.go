package user

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/model"
	"github.com/admpub/null"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

func SelectCustomer(ctx echo.Context) error {
	m := modelCustomer.NewCustomer(ctx)
	cond := &db.Compounds{}
	common.SelectPageCond(ctx, cond, `id`, `name%`)
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
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
	common.SelectPageCond(ctx, cond, `id`, `username%`)
	m := model.NewUser(ctx)
	listData := []null.StringMap{}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &listData, func(r db.Result) db.Result {
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
