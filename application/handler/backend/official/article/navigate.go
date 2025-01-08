package article

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`文章管理`),
	Action:  `official/article`,
	Icon:    `file-text-o`,
	Children: &navigate.List{
		///==========================article
		{
			Display: true,
			Name:    echo.T(`文章管理`),
			Action:  `index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加文章`),
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改文章`),
			Action:  `edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除文章`),
			Action:  `delete`,
			Icon:    `remove`,
		},
		///==========================category
		{
			Display: true,
			Name:    echo.T(`分类管理`),
			Action:  `category`,
		},
		{
			Display: false,
			Name:    echo.T(`添加分类`),
			Action:  `category_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改分类`),
			Action:  `category_edit`,
			Icon:    `pencil`,
		},
		{
			Display: false,
			Name:    echo.T(`删除分类`),
			Action:  `category_delete`,
			Icon:    `remove`,
		},
		///==========================comment
		{
			Display: true,
			Name:    echo.T(`评论管理`),
			Action:  `comment/index`,
		},
		{
			Display: false,
			Name:    echo.T(`评论管理`),
			Action:  `comment/list`,
		},
		{
			Display: false,
			Name:    echo.T(`添加评论`),
			Action:  `comment/add`,
		},
		{
			Display: false,
			Name:    echo.T(`设置评论`),
			Action:  `comment/edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除评论`),
			Action:  `comment/delete`,
		},
		///==========================friendlink
		{
			Display: true,
			Name:    echo.T(`友情链接`),
			Action:  `friendlink_index`,
		},
		{
			Display: false,
			Name:    echo.T(`添加链接`),
			Action:  `friendlink_add`,
		},
		{
			Display: false,
			Name:    echo.T(`修改链接`),
			Action:  `friendlink_edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除链接`),
			Action:  `friendlink_delete`,
		},
	},
}
