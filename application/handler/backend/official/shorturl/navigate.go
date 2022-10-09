package shorturl

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `短链接`,
	Action:  `official/short_url`,
	Icon:    `link`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `链接管理`,
			Action:  `index`,
		},
		{
			Display: false,
			Name:    `添加链接`,
			Action:  `add`,
		},
		{
			Display: false,
			Name:    `修改链接`,
			Action:  `edit/:id`,
		},
		{
			Display: false,
			Name:    `删除链接`,
			Action:  `delete/:id`,
		},
		{
			Display: true,
			Name:    `访问统计`,
			Action:  `analysis`,
		},
		{
			Display: true,
			Name:    `域名管理`,
			Action:  `domain_index`,
		},
		{
			Display: false,
			Name:    `修改域名`,
			Action:  `domain_edit/:id`,
		},
		{
			Display: false,
			Name:    `删除域名`,
			Action:  `domain_delete/:id`,
		},
	},
}
