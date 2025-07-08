package user

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/coscms/webfront/model/official"
)

func favoriteList(ctx echo.Context) error {
	customer := sessdata.Customer(ctx)
	m := official.NewCollection(ctx)
	targetType := ctx.Form(`type`, `article`)
	sort := ctx.Form(`sort`)
	sorts := make([]interface{}, 0, 1)
	switch sort {
	case `-visited`, `visited`:
		sorts = append(sorts, sort)
	case `-views`, `views`:
		sorts = append(sorts, sort)
	default:
		sorts = append(sorts, `-id`)
	}
	title := ctx.Formx(`q`).String()
	list, err := m.ListPage(targetType, customer.Id, title, sorts...)
	ctx.Set(`list`, list)
	ctx.Set(`targets`, official.CollectionTargets)
	return ctx.Render(`/user/favorite/list`, common.Err(ctx, err))
}

func favoriteDelete(ctx echo.Context) error {
	ids := param.StringSlice(ctx.FormValues(`id`)).Unique().Uint64(param.IsGreaterThanZeroElement)
	if len(ids) == 0 {
		ids = param.StringSlice(ctx.FormValues(`id[]`)).Unique().Uint64(param.IsGreaterThanZeroElement)
	}
	var err error
	var affected int64
	var customer *dbschema.OfficialCustomer
	var m *official.Collection
	var after func(isCancel ...bool) error
	if len(ids) == 0 {
		common.SendFail(ctx, ctx.T(`请选择要删除的收藏`))
		goto END
	}
	customer = sessdata.Customer(ctx)
	m = official.NewCollection(ctx)
	_, err = m.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: customer.Id},
		db.Cond{`id`: db.In(ids)},
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
	affected, err = m.Deletex(nil, db.And(
		db.Cond{`customer_id`: customer.Id},
		db.Cond{`id`: db.In(ids)},
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
		next = sessdata.URLFor(`/user/favorite/index`)
	}
	return ctx.Redirect(next)
}

func favoriteGo(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	m := official.NewCollection(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}
	customer := sessdata.Customer(ctx)
	if customer.Id != m.CustomerId {
		return nerrors.ErrUserNoPerm
	}
	if ls, ok := official.CollectionTargets[m.TargetType]; ok && ls.HasList() {
		list, err := ls.List(ctx, []*official.CollectionResponse{
			{
				OfficialCommonCollection: m.OfficialCommonCollection,
			},
		}, []uint64{m.TargetId})
		if err != nil {
			return err
		}
		var targetURL string
		if len(list) > 0 {
			targetURL = list[0].URL
		}
		if len(targetURL) > 0 {
			m.UpdateFields(nil, echo.H{
				`visited`: uint(time.Now().Unix()),
				`views`:   db.Raw(`views+1`),
			}, `id`, id)
			return ctx.Redirect(targetURL)
		}
	}
	next := ctx.Form(`next`)
	if len(next) == 0 {
		next = sessdata.URLFor(`/user/favorite/index`)
	}
	common.SendFail(ctx, ctx.T(`没有找到可以跳转的网址`))
	return ctx.Redirect(next)
}
