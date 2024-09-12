package profile

import (
	"github.com/webx-top/echo"

	mw "github.com/coscms/webfront/middleware"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func Profile(c echo.Context) error {
	profileDetail := mw.CustomerDetail(c)
	m := modelCustomer.NewCustomer(c)
	uploadConfig := m.GetUploadConfig(profileDetail.OfficialCustomer)
	var err error
	c.Set(`title`, c.T(`个人资料`))
	c.Set(`profile`, profileDetail)
	c.Set(`uploadConfig`, uploadConfig)
	return c.Render(`/user/profile/profile`, err)
}
