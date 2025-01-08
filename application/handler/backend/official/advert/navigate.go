package advert

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`广告管理`),
	Action:  `official/advert`,
	Icon:    `file-text-o`,
	Children: &navigate.List{
		///==========================item
		{
			Display: true,
			Name:    echo.T(`广告管理`),
			Action:  `index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加广告`),
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改广告`),
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除广告`),
			Action:  `delete`,
			Icon:    `remove`,
		},
		///==========================position
		{
			Display: true,
			Name:    echo.T(`广告位管理`),
			Action:  `position_index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加广告位`),
			Action:  `position_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改广告位`),
			Action:  `position_edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除广告位`),
			Action:  `position_delete`,
			Icon:    `remove`,
		},
	},
}
