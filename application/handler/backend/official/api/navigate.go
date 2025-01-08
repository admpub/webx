package api

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`外部接口`),
	Action:  `official/api`,
	Icon:    `plane`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    echo.T(`接口账号`),
			Action:  `account/index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加账号`),
			Action:  `account/add`,
		},
		{
			Display: false,
			Name:    echo.T(`修改账号`),
			Action:  `account/edit/:id`,
		},
		{
			Display: false,
			Name:    echo.T(`删除账号`),
			Action:  `account/delete/:id`,
		},
	},
}
