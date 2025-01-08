package customer

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`客户管理`),
	Action:  `official/customer`,
	Icon:    `user`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    echo.T(`客户管理`),
			Action:  `index`,
		},
		{
			Display: true,
			Name:    echo.T(`添加客户`),
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改客户`),
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除客户`),
			Action:  `delete`,
			Icon:    `remove`,
		},
		{
			Display: false,
			Name:    echo.T(`踢下线`),
			Action:  `kick`,
		},
		{
			Display: false,
			Name:    echo.T(`重新统计客户文件`),
			Action:  `recount_file`,
		},
		//level
		{
			Display: true,
			Name:    echo.T(`等级管理`),
			Action:  `level/index`,
		},
		{
			Display: true,
			Name:    echo.T(`添加等级`),
			Action:  `level/add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改等级`),
			Action:  `level/edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除等级`),
			Action:  `level/delete`,
			Icon:    `remove`,
		},
		//group_package
		{
			Display: true,
			Name:    echo.T(`等级组套餐`),
			Action:  `group_package/index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加套餐`),
			Action:  `group_package/add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改套餐`),
			Action:  `group_package/edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除套餐`),
			Action:  `group_package/delete`,
			Icon:    `remove`,
		},
		//wallet
		{
			Display: true,
			Name:    echo.T(`资产管理`),
			Action:  `wallet/index`,
		},
		{
			Display: false,
			Name:    echo.T(`修改资产`),
			Action:  `wallet/edit`,
			Icon:    `pencil`,
		},
		{
			Display: true,
			Name:    echo.T(`资产流水`),
			Action:  `wallet/flow`,
		},
		//group
		{
			Display: true,
			Name:    echo.T(`分组管理`),
			Action:  `group/index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加分组`),
			Action:  `group/add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改分组`),
			Action:  `group/edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`删除分组`),
			Action:  `group/delete`,
			Icon:    ``,
		},
		//role
		{
			Display: true,
			Name:    echo.T(`客户角色`),
			Action:  `role/index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加角色`),
			Action:  `role/add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改角色`),
			Action:  `role/edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除角色`),
			Action:  `role/delete`,
			Icon:    `remove`,
		},
		//invitation
		{
			Display: true,
			Name:    echo.T(`邀请码`),
			Action:  `invitation/index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加邀请码`),
			Action:  `invitation/add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改邀请码`),
			Action:  `invitation/edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`受邀客户列表`),
			Action:  `invitation/customer_list`,
		},
		{
			Display: false,
			Name:    echo.T(`删除邀请码`),
			Action:  `invitation/delete`,
			Icon:    `remove`,
		},
		//complaint
		{
			Display: true,
			Name:    echo.T(`客户投诉`),
			Action:  `complaint/index`,
		},
		{
			Display: false,
			Name:    echo.T(`修改投诉`),
			Action:  `complaint/edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除投诉`),
			Action:  `complaint/delete`,
			Icon:    `remove`,
		},
		{
			Display: true,
			Name:    echo.T(`充值卡管理`),
			Action:  `prepaid_card/index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加充值卡`),
			Action:  `prepaid_card/add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改充值卡`),
			Action:  `prepaid_card/edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除充值卡`),
			Action:  `prepaid_card/delete`,
			Icon:    `remove`,
		},
	},
}
