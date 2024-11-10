package user

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/coscms/webfront/model/official"
)

func favoriteList(ctx echo.Context) error {
	customer := sessdata.Customer(ctx)
	m := official.NewCollection(ctx)
	targetType := ctx.Form(`type`, `article`)
	list, err := m.ListPage(targetType, customer.Id, `-id`)
	ctx.Set(`list`, list)
	ctx.Set(`targets`, official.CollectionTargets)
	return ctx.Render(`/user/favorite/list`, common.Err(ctx, err))
}

func favoriteDelete(ctx echo.Context) error {
	ids := ctx.FormValues(`id`)
	if len(ids) == 0 {
		ids = ctx.FormValues(`id[]`)
	}
	var err error
	var affected int64
	var customer *dbschema.OfficialCustomer
	var m *official.Collection
	var nids []uint64
	var after func(isCancel ...bool) error
	if len(ids) == 0 {
		common.SendFail(ctx, ctx.T(`请选择要删除的收藏`))
		goto END
	}
	customer = sessdata.Customer(ctx)
	m = official.NewCollection(ctx)
	nids = param.Converts(ids, func(v string) uint64 {
		return param.AsUint64(v)
	})
	if len(nids) == 0 {
		common.SendFail(ctx, ctx.T(`请选择要删除的收藏`))
		goto END
	}
	_, err = m.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: customer.Id},
		db.Cond{`id`: db.In(nids)},
	))
	if err != nil {
		goto END
	}
	for _, row := range m.Objects() {
		target, ok := official.CollectionTargets[row.TargetType]
		if !ok {
			continue
		}
		after, _, err = target.Do(ctx, row.TargetId)
		if err != nil {
			continue
		}
		if after != nil {
			err = after(true)
			if err != nil {
				goto END
			}
		}
	}
	nids = param.UniqueWithFilter(nids, param.IsGreaterThanZeroElement)
	affected, err = m.Deletex(nil, `id`, db.And(
		db.Cond{`customer_id`: customer.Id},
		db.Cond{`id`: db.In(nids)},
	))

END:
	if err != nil {
		common.SendErr(ctx, err)
	} else if affected == 0 {
		common.SendFail(ctx, ctx.T(`删除失败：数据不存在`))
	} else {
		common.SendOk(ctx, ctx.T(`删除成功`))
	}

	next := ctx.Form(`next`)
	if len(next) == 0 {
		next = sessdata.URLFor(`/user/favorite`)
	}
	return ctx.Redirect(next)
}
