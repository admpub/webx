package api

import "github.com/admpub/nging/v5/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `外部接口`,
	Action:  `official/api`,
	Icon:    `plane`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `接口账号`,
			Action:  `account/index`,
		},
		{
			Display: false,
			Name:    `添加账号`,
			Action:  `account/add`,
		},
		{
			Display: false,
			Name:    `修改账号`,
			Action:  `account/edit/:id`,
		},
		{
			Display: false,
			Name:    `删除账号`,
			Action:  `account/delete/:id`,
		},
	},
}
