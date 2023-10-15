package user

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"

	dbschemaNging "github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/initialize/frontend"
	xMW "github.com/admpub/webx/application/middleware"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

// MessageUnreadCount 未读消息统计
func MessageUnreadCount(c echo.Context) error {
	customer := xMW.Customer(c)
	data := c.Data()
	m := modelCustomer.NewMessage(c)
	gids := []uint{}
	if customer.GroupId > 0 {
		gids = append(gids, customer.GroupId)
	}
	unreadUserMessages := m.CountUnread(customer.Id, gids, false)
	unreadSystemMessages := m.CountUnread(customer.Id, gids, true)
	data.SetData(echo.H{
		`user`:   unreadUserMessages,
		`system`: unreadSystemMessages,
	})
	return c.JSON(data)
}

// MessageInbox 用户消息收件箱
func MessageInbox(c echo.Context) error {
	customer := xMW.Customer(c)
	err := messageList(c, customer, false, true)
	ret := handler.Err(c, err)
	c.Set(`boxType`, `inbox`)
	c.Set(`boxTypeName`, c.T(`收件箱`))
	return c.Render(`/user/message/list`, ret)
}

// MessageOutbox 用户消息发件箱
func MessageOutbox(c echo.Context) error {
	customer := xMW.Customer(c)
	err := messageList(c, customer, false, false)
	ret := handler.Err(c, err)
	c.Set(`boxType`, `outbox`)
	c.Set(`boxTypeName`, c.T(`发件箱`))
	return c.Render(`/user/message/list`, ret)
}

func messageList(c echo.Context, customer *dbschema.OfficialCustomer, isSystemMessage bool, isRecipient bool) error {
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
	c.Request().Form().Set(`size`, `10`)
	if isRecipient {
		gids := []uint{}
		if customer.GroupId > 0 {
			gids = append(gids, customer.GroupId)
		}
		list, err = m.ListWithViewedByRecipient(customer.Id, gids, isSystemMessage, onlyUnread, cond.And())
	} else {
		list, err = m.ListWithViewedBySender(customer.Id, onlyUnread, cond.And())
	}
	if err != nil {
		return err
	}
	c.Set(`list`, list)
	return err
}

// SystemMessage 用户消息
func SystemMessage(c echo.Context) error {
	customer := xMW.Customer(c)
	err := messageList(c, customer, true, true)
	ret := handler.Err(c, err)
	c.Set(`boxType`, `system`)
	c.Set(`boxTypeName`, c.T(`系统消息`))
	return c.Render(`/user/message/list`, ret)
}

func MessageView(c echo.Context) error {
	customer := xMW.Customer(c)
	id := c.Paramx(`id`).Uint64()
	if id < 1 {
		return c.E(`id无效`)
	}
	m := modelCustomer.NewMessage(c)
	err := m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return c.E(`数据不存在`)
		}
		return err
	}
	boxType := c.Param(`type`, `inbox`)
	var (
		boxTypeName string
	)
	switch boxType {
	case `outbox`:
		boxTypeName = c.T(`发件箱`)
		if m.CustomerA != customer.Id && (customer.Uid == 0 || customer.Uid != m.UserA) {
			return c.E(`越权访问`)
		}
	case `system`:
		boxTypeName = c.T(`系统消息`)
		if !m.CheckRecvPerm(customer) {
			return c.E(`越权访问`)
		}
	default:
		boxType = `inbox`
		boxTypeName = c.T(`收件箱`)
		if !m.CheckRecvPerm(customer) {
			return c.E(`越权访问`)
		}
	}
	m.DecodeContent(m.OfficialCommonMessage)
	msgUser := m.MsgUser()
	gids := []uint{}
	if customer.GroupId > 0 {
		gids = append(gids, customer.GroupId)
	}
	err = m.View(m.OfficialCommonMessage, customer.Id, gids, `customer`) //设为已读
	if err != nil {
		return err
	}
	c.Set(`data`, m.OfficialCommonMessage)
	c.Request().Form().Set(`size`, `10`)
	_, err = common.NewLister(m.OfficialCommonMessage, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.Or(
		db.Cond{`reply_id`: m.Id},
		db.Cond{`root_id`: m.Id},
	)).Paging(c)
	rows := m.Objects()
	replyList, err := m.WithViewedByRecipient(rows, customer.Id, `customer`)
	for _, row := range replyList {
		if row.IsViewed {
			continue
		}
		if row.CustomerB > 0 && row.CustomerB == customer.Id {
			err = m.View(row.OfficialCommonMessage, customer.Id, gids, `customer`) //设为已读
			if err != nil {
				return err
			}
			continue
		}
		if row.CustomerGroupId > 0 && row.CustomerGroupId == customer.GroupId {
			err = m.View(row.OfficialCommonMessage, customer.Id, gids, `customer`) //设为已读
			if err != nil {
				return err
			}
			continue
		}
	}
	c.Set(`replyList`, replyList)
	c.Set(`msgUser`, msgUser)
	ret := handler.Err(c, err)
	c.Set(`boxType`, boxType)
	c.Set(`boxTypeName`, boxTypeName)
	return c.Render(`/user/message/view`, ret)
}

func MessageSendHandler(ctx echo.Context) error {
	m := modelCustomer.NewCustomer(ctx)
	replyID := ctx.Formx(`replyId`).Uint64()
	var customerID uint64
	if replyID > 0 {
		msgM := modelCustomer.NewMessage(ctx)
		err := msgM.Get(nil, db.Cond{`id`: replyID})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return ctx.E(`您要回复的消息不存在`)
			}
			return err
		}
		customerID = msgM.CustomerA
		if customerID < 1 {
			return ctx.E(`此消息不支持回复`)
		}
	} else {
		customerID = ctx.Formx(`customerId`).Uint64()
	}
	err := m.Get(nil, db.Cond{`id`: customerID})
	if err != nil {
		if err == db.ErrNoMoreRows {
			return ctx.E(`用户不存在`)
		}
		return err
	}
	return MessageSend(ctx, m.OfficialCustomer)
}

// MessageSend 此为公共函数，会员主页发送私信时也会调用
func MessageSend(ctx echo.Context, targetCustomer *dbschema.OfficialCustomer) error {
	replyID := ctx.Formx(`replyId`).Uint64()
	data := ctx.Data()
	customer := sessdata.Customer(ctx)
	var user *dbschemaNging.NgingUser
	if customer == nil {
		user, _ = ctx.Session().Get(`user`).(*dbschemaNging.NgingUser)
		if user == nil {
			data.SetError(common.ErrUserNotLoggedIn)
			return ctx.JSON(data)
		}
	} else {
		if replyID == 0 && targetCustomer.Id == customer.Id {
			return ctx.JSON(data.SetInfo(ctx.T(`不能给自己发私信`), 0))
		}
	}
	if ctx.IsPost() {
		data = common.VerifyCaptcha(ctx, frontend.Name, `code`)
		if common.IsFailureCode(data.GetCode()) {
			return ctx.JSON(data)
		}
		encrypted := ctx.Form(`encrypted`) == `Y`
		m := modelCustomer.NewMessage(ctx)
		m.Title = ctx.Formx(`title`).String()
		m.Content = ctx.Formx(`content`).String()
		m.ReplyId = replyID
		m.Contype = `text`
		if encrypted {
			m.Encrypted = `Y`
		} else {
			m.Encrypted = `N`
		}
		m.CustomerB = targetCustomer.Id
		if customer != nil {
			m.CustomerA = customer.Id
		} else {
			m.UserA = user.Id
		}
		_, err := m.AddData(customer, user)
		if err != nil {
			data.SetError(err)
		} else {
			data.SetInfo(ctx.T(`发送成功`))
		}
	} else {
		data.SetInfo(ctx.T(`非法操作`))
	}
	return ctx.JSON(data)
}
