package profile

import (
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/coscms/webcore/library/common"
)

// Favorites 我的收藏
// TODO: 暂不实现
func Favorites(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	var err error
	m := modelCustomer.NewFollowing(ctx)
	ctx.Request().Form().Set(`pageSize`, `20`)
	_, err = common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, db.Cond{`customer_a`: customer.Id}).Paging(ctx)
	if err != nil {
		return err
	}
	rows := m.Objects()
	list := make([]*modelCustomer.FollowingAndCustomer, len(rows))
	emptyCustomer := echo.H{}
	customerIds := []uint64{}
	for index, row := range rows {
		list[index] = &modelCustomer.FollowingAndCustomer{
			OfficialCustomerFollowing: row,
			Customer:                  emptyCustomer,
		}
		if !com.InUint64Slice(row.CustomerB, customerIds) {
			customerIds = append(customerIds, row.CustomerB)
		}
	}
	if len(customerIds) > 0 {
		customerM := modelCustomer.NewCustomer(ctx)
		_, err = customerM.ListByOffset(nil, func(r db.Result) db.Result {
			return r.Select(dbschema.DBI.OmitSelect(customerM.OfficialCustomer, modelCustomer.PrivateFieldsWithMobileEmail...)...)
		}, 0, -1, db.Cond{`id`: db.In(customerIds)})
		if err != nil {
			return err
		}
		isHTML := ctx.Format() == echo.ContentTypeHTML
		for _, customer := range customerM.Objects() {
			for index, row := range list {
				if row.CustomerB == customer.Id {
					if isHTML {
						list[index].Customer = customer.AsMap()
					} else {
						list[index].Customer = customer.AsRow().Delete(modelCustomer.PrivateFieldsWithMobileEmail...)
					}
				}
			}
		}
	}
	ctx.Set(`list`, list)
	ret := common.Err(ctx, err)
	ctx.Set(`isFollowing`, true)
	ctx.Set(`activeURL`, `/user/profile`)
	return ctx.Render(`user/profile/favorites`, ret)
}

// Likes 我的喜欢(分为喜欢的产品、新闻、评论)
// TODO: 暂不实现
func Likes(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	var err error
	m := modelCustomer.NewFollowing(ctx)
	ctx.Request().Form().Set(`pageSize`, `20`)
	_, err = common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, db.Cond{`customer_a`: customer.Id}).Paging(ctx)
	if err != nil {
		return err
	}
	rows := m.Objects()
	list := make([]*modelCustomer.FollowingAndCustomer, len(rows))
	emptyCustomer := echo.H{}
	customerIds := []uint64{}
	for index, row := range rows {
		list[index] = &modelCustomer.FollowingAndCustomer{
			OfficialCustomerFollowing: row,
			Customer:                  emptyCustomer,
		}
		if !com.InUint64Slice(row.CustomerB, customerIds) {
			customerIds = append(customerIds, row.CustomerB)
		}
	}
	if len(customerIds) > 0 {
		customerM := modelCustomer.NewCustomer(ctx)
		_, err = customerM.ListByOffset(nil, func(r db.Result) db.Result {
			return r.Select(dbschema.DBI.OmitSelect(customerM.OfficialCustomer, modelCustomer.PrivateFieldsWithMobileEmail...)...)
		}, 0, -1, db.Cond{`id`: db.In(customerIds)})
		if err != nil {
			return err
		}
		isHTML := ctx.Format() == echo.ContentTypeHTML
		for _, customer := range customerM.Objects() {
			for index, row := range list {
				if row.CustomerB == customer.Id {
					if isHTML {
						list[index].Customer = customer.AsMap()
					} else {
						list[index].Customer = customer.AsRow().Delete(modelCustomer.PrivateFieldsWithMobileEmail...)
					}
				}
			}
		}
	}
	ctx.Set(`list`, list)
	ret := common.Err(ctx, err)
	ctx.Set(`isFollowing`, true)
	ctx.Set(`activeURL`, `/user/profile`)
	return ctx.Render(`user/profile/likes`, ret)
}
