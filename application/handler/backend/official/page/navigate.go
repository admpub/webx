package page

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`页面布局`),
	Action:  `official/page`,
	Icon:    `magic`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    echo.T(`模板管理`),
			Action:  `template_index`,
			Icon:    `columns`,
		},
		{
			Display: false,
			Name:    echo.T(`修改模板`),
			Action:  `template_edit`,
			Icon:    `edit`,
		},
		{
			Display: false,
			Name:    echo.T(`切换模板`),
			Action:  `template_enable`,
			Icon:    `edit`,
		},
		{
			Display: false,
			Name:    echo.T(`修改模板配置`),
			Action:  `template_config`,
			Icon:    `edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除模板`),
			Action:  `template_delete`,
		},
		// 以下功能暂时没有用到，注释掉
		// {
		// 	Display: true,
		// 	Name:    echo.T(`页面管理`),
		// 	Action:  `index`,
		// 	Icon:    `columns`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`添加页面`),
		// 	Action:  `add`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`修改页面`),
		// 	Action:  `edit`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`删除页面`),
		// 	Action:  `delete`,
		// },
		// {
		// 	Display: true,
		// 	Name:    echo.T(`布局管理`),
		// 	Action:  `layout_index`,
		// 	Icon:    `fa-list-alt`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`添加布局`),
		// 	Action:  `layout_add`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`修改布局`),
		// 	Action:  `layout_edit`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`删除布局`),
		// 	Action:  `layout_delete`,
		// },
		// {
		// 	Display: true,
		// 	Name:    echo.T(`区块管理`),
		// 	Action:  `block_index`,
		// 	Icon:    `fa-list-alt`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`添加区块`),
		// 	Action:  `block_add`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`修改区块`),
		// 	Action:  `block_edit`,
		// },
		// {
		// 	Display: false,
		// 	Name:    echo.T(`删除区块`),
		// 	Action:  `block_delete`,
		// },
	},
}
