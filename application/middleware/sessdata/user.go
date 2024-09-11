package sessdata

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/dbschema"
)

// User 后台用户信息
func User(c echo.Context) *dbschema.NgingUser {
	user, _ := c.Session().Get(`user`).(*dbschema.NgingUser)
	return user
}
