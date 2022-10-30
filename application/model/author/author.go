package author

import (
	"fmt"

	modelNging "github.com/admpub/nging/v5/application/model"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func New(id uint64, types ...string) *Author {
	typ := `customer`
	if len(types) > 0 {
		typ = types[0]
	}
	return &Author{Type: typ, Id: id}
}

type Author struct {
	Type    string //customer/user
	Id      uint64
	Name    string
	Avatar  string
	Gender  string
	HomeURL string //用户主页网址
}

func (a *Author) Get(c echo.Context) *Author {
	if a.Id < 1 {
		a.Name = c.T(`[匿名]`)
		return a
	}
	switch a.Type {
	case `user`:
		userM := modelNging.NewUser(c)
		userM.Get(func(r db.Result) db.Result {
			return r.Select(`id`, `username`, `avatar`, `gender`)
		}, `id`, a.Id)
		a.Name = userM.Username
		a.Avatar = userM.Avatar
		a.Gender = userM.Gender
		if userM.Id < 1 {
			a.Name = c.T(`[已注销]`)
		}
	case `customer`:
		custM := modelCustomer.NewCustomer(c)
		custM.Get(func(r db.Result) db.Result {
			return r.Select(`id`, `name`, `avatar`, `gender`)
		}, `id`, a.Id)
		a.Name = custM.Name
		a.Avatar = custM.Avatar
		a.Gender = custM.Gender
		if custM.Id > 0 {
			a.HomeURL = sessdata.URLFor(`/u/` + fmt.Sprint(a.Id))
		} else {
			a.Name = c.T(`[已注销]`)
		}
	}
	return a
}
