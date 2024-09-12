package user

import (
	"github.com/coscms/webcore/library/backend"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// MessageUnreadCount 未读消息统计
func MessageUnreadCount(c echo.Context) error {
	user := backend.User(c)
	data := c.Data()
	m := modelCustomer.NewMessage(c)
	uid := uint64(user.Id)
	gids := param.String(user.RoleIds).Split(`,`).Uint()
	unreadUserMessages := m.CountUnread(uid, gids, false, `user`)
	unreadSystemMessages := m.CountUnread(uid, gids, true, `user`)
	data.SetData(echo.H{
		`user`:   unreadUserMessages,
		`system`: unreadSystemMessages,
	})
	return c.JSON(data)
}
