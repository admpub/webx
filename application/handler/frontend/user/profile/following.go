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

// IsFollowed 是否关注过
func IsFollowed(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	m := modelCustomer.NewFollowing(ctx)
	data := ctx.Data()
	m.CustomerA = customer.Id
	m.CustomerB = ctx.Formx(`uid`).Uint64()
	if m.CustomerA < 1 {
		data.SetURL(xMW.URLFor(`/sign_in`))
		data.SetError(common.ErrUserNotLoggedIn)
		return ctx.JSON(data)
	}
	if m.CustomerB < 1 {
		return ctx.JSON(data.SetInfo(ctx.T(`参数uid无效`), 0))
	}
	exists, err := m.Exists(m.CustomerA, m.CustomerB)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	if exists {
		data.SetInfo(ctx.T(`已关注`))
		return ctx.JSON(data)
	}
	data.SetInfo(ctx.T(`未关注`), 0)
	return ctx.JSON(data)
}

// Follow 关注/取关 操作
func Follow(ctx echo.Context) error {
	unfollow := ctx.Formx(`unfollow`).Bool()
	if unfollow {
		return Unfollow(ctx)
	}
	customer := xMW.Customer(ctx)
	m := modelCustomer.NewFollowing(ctx)
	data := ctx.Data()
	m.CustomerA = customer.Id
	m.CustomerB = ctx.Formx(`uid`).Uint64()
	if m.CustomerA < 1 {
		data.SetURL(xMW.URLFor(`/sign_in`))
		return ctx.JSON(data.SetInfo(ctx.T(`请先登录`), -6))
	}
	if m.CustomerB < 1 {
		return ctx.JSON(data.SetInfo(ctx.T(`参数uid无效`), 0))
	}
	if m.CustomerA == m.CustomerB {
		return ctx.JSON(data.SetInfo(ctx.T(`不能关注自己`), 0))
	}
	exists, err := m.Exists(m.CustomerA, m.CustomerB)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	if exists {
		return ctx.JSON(data.SetInfo(ctx.T(`您已经关注过了`), 0))
	}
	_, err = m.Add()
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	data.SetInfo(ctx.T(`关注成功`))
	return ctx.JSON(data)
}

// Unfollow 取消关注 操作
func Unfollow(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	m := modelCustomer.NewFollowing(ctx)
	data := ctx.Data()
	m.CustomerA = customer.Id
	m.CustomerB = ctx.Formx(`uid`).Uint64()
	if m.CustomerA < 1 {
		data.SetURL(xMW.URLFor(`/sign_in`))
		return ctx.JSON(data.SetInfo(ctx.T(`请先登录`), -6))
	}
	if m.CustomerB < 1 {
		return ctx.JSON(data.SetInfo(ctx.T(`参数uid无效`), 0))
	}
	exists, err := m.Exists(m.CustomerA, m.CustomerB)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	if exists {
		err = m.Delete(nil, db.And(
			db.Cond{`customer_a`: m.CustomerA},
			db.Cond{`customer_b`: m.CustomerB},
		))
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		data.SetInfo(ctx.T(`取关成功`))
		return ctx.JSON(data)
	}
	data.SetInfo(ctx.T(`您并没有关注过该会员`), 0)
	return ctx.JSON(data)
}

// Following 我关注的
func Following(ctx echo.Context) error {
	if ctx.IsPost() {
		return Follow(ctx)
	}
	customer := xMW.Customer(ctx)
	err := FollowingBy(ctx, customer.Id)
	ctx.Set(`activeURL`, `/user/profile`)
	ctx.Set(`title`, ctx.T(`我关注的用户`))
	return ctx.Render(`user/profile/following`, common.Err(ctx, err))
}

func FollowingBy(ctx echo.Context, customerID uint64) error {
	m := modelCustomer.NewFollowing(ctx)
	ctx.Request().Form().Set(`pageSize`, `20`)
	_, err := common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, db.Cond{`customer_a`: customerID}).Paging(ctx)
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
	ctx.Set(`isFollowing`, true)
	return err
}

// Followers 关注我的
func Followers(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	err := FollowersBy(ctx, customer.Id)
	ctx.Set(`activeURL`, `/user/profile`)
	ctx.Set(`title`, ctx.T(`关注我的用户`))
	return ctx.Render(`user/profile/following`, common.Err(ctx, err))
}

func FollowersBy(ctx echo.Context, customerID uint64) error {
	m := modelCustomer.NewFollowing(ctx)
	ctx.Request().Form().Set(`pageSize`, `20`)
	_, err := common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, db.Cond{`customer_b`: customerID}).Paging(ctx)
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
		if !com.InUint64Slice(row.CustomerA, customerIds) {
			customerIds = append(customerIds, row.CustomerA)
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
				if row.CustomerA == customer.Id {
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
	ctx.Set(`isFollowing`, false)
	return err
}
