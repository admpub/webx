package agent

import "github.com/admpub/nging/v5/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `代理商`,
	Action:  `official/agent`,
	Icon:    `trophy`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `代理商列表`,
			Action:  `index`,
		},
		{
			Display: false,
			Name:    `添加代理商`,
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改代理商`,
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    `删除代理商`,
			Action:  `delete`,
			Icon:    `remove`,
		},
		//product
		{
			Display: true,
			Name:    `代理产品`,
			Action:  `product_index`,
		},
		{
			Display: false,
			Name:    `添加代理产品`,
			Action:  `product_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改代理产品`,
			Action:  `product_edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    `删除代理产品`,
			Action:  `product_delete`,
			Icon:    `remove`,
		},
		//level
		{
			Display: true,
			Name:    `代理等级`,
			Action:  `level_index`,
		},
		{
			Display: false,
			Name:    `添加代理等级`,
			Action:  `level_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改代理等级`,
			Action:  `level_edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    `删除代理等级`,
			Action:  `level_delete`,
			Icon:    `remove`,
		},
	},
}
