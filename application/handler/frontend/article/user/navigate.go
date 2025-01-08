package user

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webfront/initialize/frontend/usernav"
	"github.com/webx-top/echo"
)

func init() {
	usernav.LeftNavigate.Add(-1, &navigate.Item{
		Display: true,
		Name:    echo.T(`我的文章`),
		Action:  `article`,
		Icon:    `file-text-o`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     echo.T(`我的文章`),
				Action:   `list`,
				Icon:     `file-text-o`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     echo.T(`投稿`),
				Action:   `create`,
				Icon:     `plus`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     echo.T(`修改文章`),
				Action:   `edit/:id`,
				Icon:     `pencil`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     echo.T(`删除文章`),
				Action:   `delete/:id`,
				Icon:     `times`,
				Children: &navigate.List{},
			},
		},
	})
}
