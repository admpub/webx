package customer

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

type FollowingAndCustomer struct {
	*dbschema.OfficialCustomerFollowing
	Customer echo.H
}
