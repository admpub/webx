package user

import (
	"github.com/admpub/nging/v4/application/registry/navigate"
	userNav "github.com/admpub/webx/application/handler/frontend/user/navigate"
)

func init() {
	userNav.LeftNavigate.Add(-1, &navigate.Item{
		Display: true,
		Name:    `我的文章`,
		Action:  `article`,
		Icon:    `file-text-o`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     `我的文章`,
				Action:   `list`,
				Icon:     `table`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `投稿`,
				Action:   `create`,
				Icon:     `plus`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `修改文章`,
				Action:   `edit/:id`,
				Icon:     `pencil`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `删除文章`,
				Action:   `delete/:id`,
				Icon:     `times`,
				Children: &navigate.List{},
			},
		},
	})
}
