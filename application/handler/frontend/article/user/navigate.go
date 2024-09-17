package user

import (
	"github.com/coscms/webcore/registry/navigate"
	"github.com/coscms/webfront/initialize/frontend/usernav"
)

func init() {
	usernav.LeftNavigate.Add(-1, &navigate.Item{
		Display: true,
		Name:    `我的文章`,
		Action:  `article`,
		Icon:    `file-text-o`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     `我的文章`,
				Action:   `list`,
				Icon:     `file-text-o`,
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
