package shorturl

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`短链接`),
	Action:  `official/short_url`,
	Icon:    `link`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    echo.T(`链接管理`),
			Action:  `index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加链接`),
			Action:  `add`,
		},
		{
			Display: false,
			Name:    echo.T(`修改链接`),
			Action:  `edit/:id`,
		},
		{
			Display: false,
			Name:    echo.T(`删除链接`),
			Action:  `delete/:id`,
		},
		{
			Display: true,
			Name:    echo.T(`访问统计`),
			Action:  `analysis`,
		},
		{
			Display: true,
			Name:    echo.T(`域名管理`),
			Action:  `domain_index`,
		},
		{
			Display: false,
			Name:    echo.T(`修改域名`),
			Action:  `domain_edit/:id`,
		},
		{
			Display: false,
			Name:    echo.T(`删除域名`),
			Action:  `domain_delete/:id`,
		},
	},
}
