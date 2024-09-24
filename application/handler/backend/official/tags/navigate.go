package tags

import "github.com/coscms/webcore/library/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `标签管理`,
	Action:  `official/tags`,
	Icon:    `tags`,
	Children: &navigate.List{
		///==========================article
		{
			Display: true,
			Name:    `标签管理`,
			Action:  `index`,
		},
		/* 标签为自动添加
		{
			Display: true,
			Name:          `添加标签`,
			Action:        `add`,
			Icon:          `plus`,
		},
		*/
		{
			Display: false,
			Name:    `修改标签`,
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    `删除标签`,
			Action:  `delete`,
			Icon:    `remove`,
		},
	},
}
