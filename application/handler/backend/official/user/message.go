package user

import (
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/model"
	xschema "github.com/coscms/webfront/dbschema"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

// MessageInbox 用户消息收件箱
func MessageInbox(c echo.Context) error {
	user := backend.User(c)
	err := messageList(c, user, false, true)
	ret := common.Err(c, err)
	c.Set(`boxType`, `inbox`)
	c.Set(`boxTypeName`, c.T(`收件箱`))
	return c.Render(`official/user/message/list`, ret)
}

// MessageOutbox 用户消息发件箱
func MessageOutbox(c echo.Context) error {
	user := backend.User(c)
	err := messageList(c, user, false, false)
	ret := common.Err(c, err)
	c.Set(`boxType`, `outbox`)
	c.Set(`boxTypeName`, c.T(`发件箱`))
	return c.Render(`official/user/message/list`, ret)
}

func messageList(c echo.Context, user *dbschema.NgingUser, isSystemMessage bool, isRecipient bool) error {
	m := modelCustomer.NewMessage(c)
	typ := c.Form(`type`)
	q := c.Formx(`q`).String()
	var onlyUnread bool
	if typ == `unread` {
		onlyUnread = true
	}
	var cond db.Compound
	if len(q) > 0 {
		cond = db.And(
			db.Cond{`encrypted`: `N`},
			db.Or(
				db.Cond{`title`: db.Like(`%` + q + `%`)},
				db.Cond{`content`: db.Like(`%` + q + `%`)},
			),
		)
	}
	var (
		list []*modelCustomer.MessageWithViewed
		err  error
	)
	c.Request().Form().Set(`size`, `10`)
	uid := uint64(user.Id)
	gids := param.String(user.RoleIds).Split(`,`).Uint()
	if isRecipient {
		list, err = m.ListWithViewedByRecipient(uid, gids, isSystemMessage, onlyUnread, cond, `user`)
	} else {
		list, err = m.ListWithViewedBySender(uid, onlyUnread, cond, `user`)
	}
	if err != nil {
		return err
	}
	c.Set(`list`, list)
	return err
}

// SystemMessage 用户消息
func SystemMessage(c echo.Context) error {
	user := backend.User(c)
	err := messageList(c, user, true, true)
	ret := common.Err(c, err)
	c.Set(`boxType`, `system`)
	c.Set(`boxTypeName`, c.T(`系统消息`))
	return c.Render(`official/user/message/list`, ret)
}

func MessageView(c echo.Context) error {
	user := backend.User(c)
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
	var boxTypeName string
	switch boxType {
	case `outbox`:
		boxTypeName = c.T(`发件箱`)
		if m.UserA > 0 && m.UserA != user.Id {
			return c.E(`越权访问`)
		}
	case `system`:
		boxTypeName = c.T(`系统消息`)
		if m.UserB > 0 && m.UserB != user.Id {
			return c.E(`越权访问`)
		}
	default:
		boxType = `inbox`
		boxTypeName = c.T(`收件箱`)
		if m.UserB > 0 && m.UserB != user.Id {
			return c.E(`越权访问`)
		}
	}
	m.DecodeContent(m.OfficialCommonMessage)
	msgUser := m.MsgUser()
	uid := uint64(user.Id)
	gids := param.String(user.RoleIds).Split(`,`).Uint()
	err = m.View(m.OfficialCommonMessage, uid, gids, `user`) //设为已读
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
	if err != nil {
		return err
	}
	rows := m.Objects()
	replyList, err := m.WithViewedByRecipient(rows, uid, `user`)
	for _, row := range replyList {
		if row.IsViewed {
			continue
		}
		if row.UserB > 0 && row.UserB == user.Id {
			err = m.View(row.OfficialCommonMessage, uid, gids, `user`) //设为已读
			if err != nil {
				return err
			}
			continue
		}
		if row.UserRoleId > 0 && com.InUintSlice(row.CustomerGroupId, gids) {
			err = m.View(row.OfficialCommonMessage, uid, gids, `user`) //设为已读
			if err != nil {
				return err
			}
			continue
		}
	}
	c.Set(`replyList`, replyList)
	c.Set(`msgUser`, msgUser)
	ret := common.Err(c, err)
	c.Set(`boxType`, boxType)
	c.Set(`boxTypeName`, boxTypeName)
	return c.Render(`official/user/message/view`, ret)
}

func MessageSendHandler(ctx echo.Context) error {
	if len(ctx.Form(`operate`)) > 0 {
		return SelectOwner(ctx, ctx.Form(`operate`))
	}
	var (
		err         error
		userID      uint
		customerID  uint64
		replySender string
		msgM        = modelCustomer.NewMessage(ctx)
		userM       = model.NewUser(ctx)
		custM       = modelCustomer.NewCustomer(ctx)
	)
	replyID := ctx.Formx(`replyId`).Uint64()
	if replyID > 0 {
		err := msgM.Get(nil, db.Cond{`id`: replyID})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return ctx.E(`您要回复的消息不存在`)
			}
			return err
		}
		userID = msgM.UserA
		customerID = msgM.CustomerA
		if userID < 1 && customerID < 1 {
			return ctx.E(`此消息不支持回复`)
		}
		if customerID > 0 {
			err = custM.Get(nil, db.Cond{`id`: customerID})
			replySender = custM.OfficialCustomer.Name
		} else {
			err = userM.Get(nil, db.Cond{`id`: userID})
			replySender = userM.NgingUser.Username
		}
		if err != nil {
			if err == db.ErrNoMoreRows {
				if customerID > 0 {
					return ctx.E(`客户不存在`)
				} else {
					return ctx.E(`用户不存在`)
				}
			}
			return err
		}
	}
	if ctx.IsPost() {
		if replyID == 0 {
			userID = ctx.Formx(`userId`).Uint()
			customerID = ctx.Formx(`customerId`).Uint64()
			if userID == 0 && customerID == 0 {
				recipientType := ctx.Form(`recipientType`)
				if recipientType == `customer` {
					customerID = ctx.Formx(`recipientId`).Uint64()
				} else {
					userID = ctx.Formx(`recipientId`).Uint()
				}
			}
			if userID > 0 {
				err = userM.Get(nil, db.Cond{`id`: userID})
			} else if customerID > 0 {
				err = custM.Get(nil, db.Cond{`id`: customerID})
			}
			if err != nil {
				if err == db.ErrNoMoreRows {
					if userID > 0 {
						return ctx.E(`用户不存在`)
					} else {
						return ctx.E(`客户不存在`)
					}
				}
				return err
			}
		}
		return MessageSend(ctx, userM.NgingUser, custM.OfficialCustomer)
	}
	ctx.Set(`replyId`, replyID)
	ctx.Set(`replyMsg`, msgM.OfficialCommonMessage)
	ctx.Set(`replySender`, replySender)
	ctx.Set(`boxType`, `send`)
	ctx.Set(`boxTypeName`, ctx.T(`发送消息`))
	return ctx.Render(`official/user/message/send`, common.Err(ctx, err))
}

// MessageSend 发送消息
func MessageSend(ctx echo.Context, targetUser *dbschema.NgingUser, targetCustomer *xschema.OfficialCustomer) error {
	replyID := ctx.Formx(`replyId`).Uint64()
	data := ctx.Data()
	user := backend.User(ctx)
	if user == nil {
		return ctx.JSON(data.SetError(common.ErrUserNotLoggedIn))
	}
	if replyID == 0 {
		if targetUser.Id > 0 && targetUser.Id == user.Id {
			return ctx.JSON(data.SetInfo(ctx.T(`不能给自己发私信`), 0))
		}
		if targetCustomer.Id > 0 && targetCustomer.Uid == user.Id {
			return ctx.JSON(data.SetInfo(ctx.T(`不能给自己发私信`), 0))
		}
	}
	if ctx.IsPost() {
		// data = common.VerifyCaptcha(ctx, frontend.Name, `code`)
		// if common.IsFailureCode(data.GetCode()) {
		// 	return ctx.JSON(data)
		// }
		encrypted := ctx.Form(`encrypted`) == `Y`
		m := modelCustomer.NewMessage(ctx)
		m.Title = ctx.Formx(`title`).String()
		m.Content = ctx.Formx(`content`).String()
		m.ReplyId = replyID
		m.Contype = ctx.Formx(`contentType`, `text`).String()
		if encrypted {
			m.Encrypted = `Y`
		} else {
			m.Encrypted = `N`
		}
		if targetUser.Id > 0 {
			m.UserB = targetUser.Id
		} else {
			m.UserB = targetCustomer.Uid // 管理员给后台管理员X的前台账号发送的消息，允许管理员X在后台查看
			m.CustomerB = targetCustomer.Id
		}
		if m.CustomerB < 1 {
			return ctx.JSON(data.SetError(ctx.NewError(code.InvalidParameter, `请选择收信人`).SetZone(`customerId`)))
		}
		m.UserA = user.Id
		_, err := m.AddData(nil, user)
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

// MessageDelete 删除消息
func MessageDelete(ctx echo.Context) error {
	ids := param.StringSlice(ctx.FormValues(`messageId[]`)).Uint64()
	data := ctx.Data()
	user := backend.User(ctx)
	m := modelCustomer.NewMessage(ctx)
	if ctx.IsPost() {
		err := m.Delete(db.And(
			db.Cond{`id`: db.In(ids)},
			db.Or(
				db.Cond{`user_a`: user.Id},
				db.Cond{`user_b`: user.Id},
			),
		))
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		data.SetInfo(ctx.T(`删除成功`))
	} else {
		data.SetInfo(ctx.T(`非法操作`), 0)
	}
	return ctx.JSON(data)
}
