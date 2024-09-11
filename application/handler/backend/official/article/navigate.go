package article

import "github.com/coscms/webcore/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `文章管理`,
	Action:  `official/article`,
	Icon:    `file-text-o`,
	Children: &navigate.List{
		///==========================article
		{
			Display: true,
			Name:    `文章管理`,
			Action:  `index`,
		},
		{
			Display: false,
			Name:    `添加文章`,
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改文章`,
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    `删除文章`,
			Action:  `delete`,
			Icon:    `remove`,
		},
		///==========================category
		{
			Display: true,
			Name:    `分类管理`,
			Action:  `category`,
		},
		{
			Display: false,
			Name:    `添加分类`,
			Action:  `category_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改分类`,
			Action:  `category_edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    `删除分类`,
			Action:  `category_delete`,
			Icon:    `remove`,
		},
		///==========================comment
		{
			Display: true,
			Name:    `评论管理`,
			Action:  `comment/index`,
		},
		{
			Display: false,
			Name:    `评论管理`,
			Action:  `comment/list`,
		},
		{
			Display: false,
			Name:    `添加评论`,
			Action:  `comment/add`,
		},
		{
			Display: false,
			Name:    `设置评论`,
			Action:  `comment/edit`,
		},
		{
			Display: false,
			Name:    `删除评论`,
			Action:  `comment/delete`,
		},
		///==========================friendlink
		{
			Display: true,
			Name:    `友情链接`,
			Action:  `friendlink_index`,
		},
		{
			Display: false,
			Name:    `添加链接`,
			Action:  `friendlink_add`,
		},
		{
			Display: false,
			Name:    `修改链接`,
			Action:  `friendlink_edit`,
		},
		{
			Display: false,
			Name:    `删除链接`,
			Action:  `friendlink_delete`,
		},
	},
}
