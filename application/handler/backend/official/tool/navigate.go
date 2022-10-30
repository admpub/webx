package tool

import "github.com/admpub/nging/v5/application/registry/navigate"

var TopNavigate = navigate.List{
	&navigate.Item{
		Display: true,
		Name:    `地区管理`,
		Action:  `area/index`,
	},
	&navigate.Item{
		Display: false,
		Name:    `添加管理`,
		Action:  `area/add`,
	},
	&navigate.Item{
		Display: false,
		Name:    `修改地区`,
		Action:  `area/edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    `删除地区`,
		Action:  `area/delete`,
	},
	&navigate.Item{
		Display: true,
		Name:    `中文分词`,
		Action:  `segment`,
	},
	&navigate.Item{
		Display: true,
		Name:    `敏感词`,
		Action:  `sensitive/index`,
	},
	&navigate.Item{
		Display: false,
		Name:    `添加敏感词`,
		Action:  `sensitive/add`,
		Icon:    `pencil`,
	},
	&navigate.Item{
		Display: false,
		Name:    `修改敏感词`,
		Action:  `sensitive/edit`,
		Icon:    `pencil`,
	},
	&navigate.Item{
		Display: false,
		Name:    `删除敏感词`,
		Action:  `sensitive/delete`,
		Icon:    `remove`,
	},
}
