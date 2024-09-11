package navigate

import (
	"github.com/coscms/webcore/registry/navigate"
)

var TopNavigate = &navigate.List{
	{
		Display: true,
		Name:    `个人消息`,
		Action:  `message`,
		Icon:    `email`,
		Children: &navigate.List{
			{
				Display:  true,
				Name:     `收件箱`,
				Action:   `inbox`,
				Icon:     `email`,
				Children: &navigate.List{},
			},
			{
				Display:  false,
				Name:     `发件箱`,
				Action:   `outbox`,
				Children: &navigate.List{},
			},
			{
				Display:  true,
				Name:     `系统消息`,
				Action:   `system`,
				Icon:     `bell`,
				Children: &navigate.List{},
			},
			{
				Display:   false,
				Name:      `查看消息内容`,
				Action:    `view/:type/:id`,
				Unlimited: true,
				Children:  &navigate.List{},
			},
			{
				Display:   false,
				Name:      `未读消息统计`,
				Action:    `unread_count`,
				Unlimited: true,
				Children:  &navigate.List{},
			},
			{
				Display:  false,
				Name:     `发送私信`,
				Action:   `send`,
				Children: &navigate.List{},
			},
		},
	},
}
