package profile

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func Settings(ctx echo.Context) error {
	customer := xMW.Customer(ctx)
	var err error
	if ctx.IsPost() {
		m := modelCustomer.NewCustomer(ctx)
		cond := db.Cond{`id`: customer.Id}
		err = m.Get(nil, cond)
		if err != nil {
			goto END
		}
		password := ctx.Formx(`password`).String()
		if len(password) > 0 {
			if password != ctx.Form(`password2`) {
				err = ctx.E(`登录密码不一致`)
				goto END
			}
			m.Password = password
		}
		safePwd := ctx.Formx(`safePwd`).String()
		if len(safePwd) > 0 {
			if safePwd != ctx.Form(`safePwd2`) {
				err = ctx.E(`安全密码不一致`)
				goto END
			}
			m.SafePwd = safePwd
		}
		m.RealName = ctx.Formx(`realName`).String()
		m.IdCardNo = ctx.Formx(`idCardNo`).String()
		if m.EmailBind == `N` {
			m.Email = ctx.Formx(`email`).String()
		}
		if m.MobileBind == `N` {
			m.Mobile = ctx.Formx(`mobile`).String()
		}
		m.Gender = ctx.Formx(`gender`).String()
		m.Description = ctx.Formx(`description`).String()
		m.Avatar = ctx.Formx(`avatar`).String()
		err = m.Edit(nil, cond)
		if err == nil {
			m.SetSession()
			return ctx.Redirect(xMW.URLFor(`/user/profile/settings`))
		}
	}

END:
	ret := handler.Err(ctx, err)
	ctx.Set(`activeURL`, `/user/profile`)
	ctx.Set(`title`, ctx.T(`账号设置`))
	return ctx.Render(`user/profile/settings`, ret)
}
