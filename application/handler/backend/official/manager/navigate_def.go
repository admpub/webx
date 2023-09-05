package manager

import "github.com/admpub/nging/v5/application/registry/navigate"

var TopNavigate = navigate.List{
	&navigate.Item{
		Display: true,
		Name:    `消息管理`,
		Action:  `message/index`,
	},
	&navigate.Item{
		Display: false,
		Name:    `删除消息`,
		Action:  `message/delete`,
	},
	&navigate.Item{
		Display: false,
		Name:    `查看消息`,
		Action:  `message/view/:id`,
	},
	&navigate.Item{
		Display: true,
		Name:    `菜单管理`,
		Action:  `navigate/index`,
	},
	&navigate.Item{
		Display: false,
		Name:    `添加菜单`,
		Action:  `navigate/add`,
	},
	&navigate.Item{
		Display: false,
		Name:    `修改菜单`,
		Action:  `navigate/edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    `删除菜单`,
		Action:  `navigate/delete`,
	},
	&navigate.Item{
		Display: true,
		Name:    `重启前台`,
		Action:  `frontend/reboot`,
		Target:  `ajax`,
	},
	&navigate.Item{
		Display: true,
		Name:    `自定义页面`,
		Action:  `frontend/route_page`,
	},
	&navigate.Item{
		Display: false,
		Name:    `添加自定义页面`,
		Action:  `frontend/route_page_add`,
	},
	&navigate.Item{
		Display: false,
		Name:    `修改自定义页面`,
		Action:  `frontend/route_page_edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    `删除自定义页面`,
		Action:  `frontend/route_page_delete`,
	},
	&navigate.Item{
		Display: true,
		Name:    `自定义网址`,
		Action:  `frontend/route_rewrite`,
	},
	&navigate.Item{
		Display: false,
		Name:    `添加自定义网址`,
		Action:  `frontend/route_rewrite_add`,
	},
	&navigate.Item{
		Display: false,
		Name:    `修改自定义网址`,
		Action:  `frontend/route_rewrite_edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    `删除自定义网址`,
		Action:  `frontend/route_rewrite_delete`,
	},
}
