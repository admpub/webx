package manager

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

// MessageIndex 所有消息列表
func MessageIndex(c echo.Context) error {
	err := messageList(c)
	ret := common.Err(c, err)
	return c.Render(`official/manager/message/index`, ret)
}

func messageList(c echo.Context) error {
	m := modelCustomer.NewMessage(c)
	typ := c.Form(`type`)
	q := c.Formx(`q`).String()
	var onlyUnread bool
	if typ == `unread` {
		onlyUnread = true
	}
	cond := db.NewCompounds()
	if len(q) > 0 {
		cond.AddKV(`encrypted`, `N`)
		cond.From(mysql.SearchField(`~title+content`, q))
	}
	var (
		list []*modelCustomer.MessageWithViewed
		err  error
	)
	//c.Request().Form().Set(`size`, `20`)
	list, err = m.ListAll(onlyUnread, cond.And())
	if err != nil {
		return err
	}
	c.Set(`listData`, list)
	return err
}

func MessageView(c echo.Context) error {
	id := c.Paramx(`id`).Uint64()
	if id < 1 {
		return c.NewError(code.InvalidParameter, `id无效`).SetZone(`id`)
	}
	m := modelCustomer.NewMessage(c)
	err := m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return c.NewError(code.DataNotFound, `数据不存在`)
		}
		return err
	}
	m.DecodeContent(m.OfficialCommonMessage)
	data, err := m.GetWithViewed(m.OfficialCommonMessage)
	c.Set(`data`, data)
	ret := common.Err(c, err)
	return c.Render(`official/manager/message/view`, ret)
}

// MessageDelete 删除消息
func MessageDelete(ctx echo.Context) error {
	var ids []uint64
	id := ctx.Formx(`id`).Uint64()
	if id > 0 {
		ids = append(ids, id)
	} else {
		ids = param.StringSlice(ctx.FormValues(`messageId[]`)).Uint64()
	}
	data := ctx.Data()
	m := modelCustomer.NewMessage(ctx)
	if len(ids) > 0 {
		err := m.Delete(db.And(
			db.Cond{`id`: db.In(ids)},
		))
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		data.SetInfo(ctx.T(`删除成功`))
	} else {
		data.SetInfo(ctx.T(`非法操作`), 0)
	}
	if ctx.Format() == `html` {
		if data.GetCode().Int() == 1 {
			common.SendOk(ctx, param.AsString(data.GetInfo()))
		} else {
			common.SendFail(ctx, param.AsString(data.GetInfo()))
		}
	}
	return ctx.Redirect(backend.URLFor(`/manager/message/index`))
}
