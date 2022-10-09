package profile

import (
	mw "github.com/admpub/webx/application/middleware"
	"github.com/webx-top/echo"
)

func Profile(c echo.Context) error {
	var err error
	c.Set(`title`, c.T(`个人资料`))
	c.Set(`profile`, mw.CustomerDetail(c))
	return c.Render(`/user/profile/profile`, err)
}
