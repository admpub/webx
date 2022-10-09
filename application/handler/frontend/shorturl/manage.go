package shorturl

import (
	"fmt"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelShorturl "github.com/admpub/webx/application/model/official/shorturl"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

// List 用户短网址列表
func List(ctx echo.Context) error {
	customer := sessdata.Customer(ctx)
	var err error
	m := modelShorturl.NewShortURL(ctx)
	cond := db.Compounds{
		db.Cond{`owner_id`: customer.Id},
		db.Cond{`owner_type`: `customer`},
	}
	shortID := ctx.Formx(`q`).String()
	if len(shortID) > 0 {
		cond.AddKV(`short_url`, shortID)
	}
	sorts := common.Sorts(ctx, `official_short_url`, `-id`)
	_, err = common.NewLister(m.URL, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	ctx.Set(`list`, m.URL.Objects())
	return ctx.Render(`short_url/list`, handler.Err(ctx, err))
}

// Create 创建短网址
func Create(ctx echo.Context) error {
	customer := sessdata.Customer(ctx)
	var err error
	m := modelShorturl.NewShortURL(ctx)
	if ctx.IsPost() {
		m.URL.OwnerId = customer.Id
		m.URL.OwnerType = `customer`
		m.URL.LongUrl = ctx.Formx(`url`).String()
		m.URL.Password = ctx.Formx(`password`).String()
		_, err = m.Add()
		if err != nil {
			goto END
		}
		return ctx.Redirect(sessdata.URLFor(`/user/short_url/list`))
	}

END:
	ctx.Set(`activeURL`, `/user/short_url/list`)
	ctx.Set(`title`, ctx.T(`添加短链接`))
	return ctx.Render(`short_url/edit`, handler.Err(ctx, err))
}

// Edit 修改短网址
func Edit(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数“%s”值无效`, `id`)
	}
	customer := sessdata.Customer(ctx)
	m := modelShorturl.NewShortURL(ctx)
	err := m.URL.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `短网址不存在`)
		}
		return err
	}
	if m.URL.OwnerType != `customer` || m.URL.OwnerId != customer.Id {
		return ctx.NewError(code.NonPrivileged, `越权操作！您没有权限修改此数据`)
	}
	if ctx.IsPost() {
		m.URL.LongUrl = ctx.Formx(`url`).String()
		m.URL.Password = ctx.Formx(`password`).String()
		err = m.Edit(nil, `id`, m.URL.Id)
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(sessdata.URLFor(`/user/short_url/edit/` + fmt.Sprint(id)))
	}
	echo.StructToForm(ctx, m.URL, ``, echo.LowerCaseFirstLetter)
	ctx.Request().Form().Set(`url`, ctx.Form(`longUrl`))

END:
	ctx.Set(`activeURL`, `/user/short_url/list`)
	ctx.Set(`title`, ctx.T(`修改短链接`))
	return ctx.Render(`short_url/edit`, common.Err(ctx, err))
}

// Delete 用户删除短网址
func Delete(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数“%s”值无效`, `id`)
	}
	customer := sessdata.Customer(ctx)
	m := modelShorturl.NewShortURL(ctx)
	err := m.URL.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `短网址不存在`)
		}
		return err
	}
	if m.URL.OwnerType != `customer` || m.URL.OwnerId != customer.Id {
		return ctx.NewError(code.NonPrivileged, `越权操作！您没有权限删除此数据`)
	}
	return ctx.Redirect(sessdata.URLFor(`/user/short_url/list`))
}

// Analysis 用户短网址访问统计
func Analysis(ctx echo.Context) error {
	customer := sessdata.Customer(ctx)
	m := modelShorturl.NewShortURL(ctx)
	_, err := m.Visit.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, 0, 50, db.And(
		db.Cond{`owner_id`: customer.Id},
		db.Cond{`owner_type`: `customer`},
	))
	if err != nil {
		return err
	}
	ctx.Set(`lasts`, m.Visit.Objects())
	ctx.Set(`activeURL`, `/user/short_url/list`)
	return ctx.Render(`short_url/analysis`, handler.Err(ctx, err))
}
