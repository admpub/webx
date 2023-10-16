package agent

import (
	"github.com/admpub/nging/v5/application/registry/navigate"
)

var LeftNavigate = &navigate.List{
	{
		Display: true,
		Name:    `代理`,
		Action:  `agent`,
		Icon:    `link`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     `我的代理`,
				Action:   `index`,
				Icon:     `table`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `修改代理资料`,
				Action:   `edit`,
				Icon:     `pencil`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `我的邀请`,
				Action:   `invited`,
				Icon:     `table`,
				Children: &navigate.List{},
			},
		},
	},
}
