package tags

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`标签管理`),
	Action:  `official/tags`,
	Icon:    `tags`,
	Children: &navigate.List{
		///==========================article
		{
			Display: true,
			Name:    echo.T(`标签管理`),
			Action:  `index`,
		},
		{
			Display: true,
			Name:    echo.T(`添加标签`),
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改标签`),
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除标签`),
			Action:  `delete`,
			Icon:    `remove`,
		},
	},
}
