package sessdata

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/dbschema"
)

// User 后台用户信息
func User(c echo.Context) *dbschema.NgingUser {
	user, _ := c.Get(`user`).(*dbschema.NgingUser)
	return user
}
