package manager

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var TopNavigate = navigate.List{
	&navigate.Item{
		Display: true,
		Name:    echo.T(`消息管理`),
		Action:  `message/index`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`删除消息`),
		Action:  `message/delete`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`查看消息`),
		Action:  `message/view/:id`,
	},
	&navigate.Item{
		Display: true,
		Name:    echo.T(`前台菜单`),
		Action:  `navigate/index`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`添加前台菜单`),
		Action:  `navigate/add`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`修改前台菜单`),
		Action:  `navigate/edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`删除前台菜单`),
		Action:  `navigate/delete`,
	},
	&navigate.Item{
		Display: true,
		Name:    echo.T(`重启前台`),
		Action:  `frontend/reboot`,
		Target:  `ajax`,
	},
	&navigate.Item{
		Display: true,
		Name:    echo.T(`自定义页面`),
		Action:  `frontend/route_page`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`添加自定义页面`),
		Action:  `frontend/route_page_add`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`修改自定义页面`),
		Action:  `frontend/route_page_edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`删除自定义页面`),
		Action:  `frontend/route_page_delete`,
	},
	&navigate.Item{
		Display: true,
		Name:    echo.T(`自定义网址`),
		Action:  `frontend/route_rewrite`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`添加自定义网址`),
		Action:  `frontend/route_rewrite_add`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`修改自定义网址`),
		Action:  `frontend/route_rewrite_edit`,
	},
	&navigate.Item{
		Display: false,
		Name:    echo.T(`删除自定义网址`),
		Action:  `frontend/route_rewrite_delete`,
	},
}
