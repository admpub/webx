package user

import (
	"github.com/webx-top/echo"
)

// 个人资料

func Index(c echo.Context) error {
	return c.Redirect(`/user/profile/detail`)
	var err error
	return c.Render(`/user/index`, err)
}
