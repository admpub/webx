package transformCustomer

import "github.com/webx-top/echo/param"

var Detail = map[string]param.Transfer{
	`Group`: param.Tf(`group`, nil),
	`Level`: param.Tf(`level`, func(value interface{}, row param.Store) interface{} {
		mp, ok := value.(param.Store)
		if !ok {
			return value
		}
		return param.Store{
			`id`:          mp.Get(`Id`),
			`short`:       mp.Get(`Short`),
			`name`:        mp.Get(`Name`),
			`description`: mp.Get(`Description`),
			`icon_image`:  mp.Get(`IconImage`),
			`icon_class`:  mp.Get(`IconClass`),
			`color`:       mp.Get(`Color`),
			`bgcolor`:     mp.Get(`Bgcolor`),
		}
	}),
	`Agent`: param.Tf(`agent`, func(value interface{}, row param.Store) interface{} {
		mp, ok := value.(param.Store)
		if !ok {
			return value
		}
		return param.Store{
			`id`:          mp.Get(`Id`),
			`name`:        mp.Get(`Name`),
			`description`: mp.Get(`Description`),
		}
	}),
	`Roles.Id`:          param.Tf(`roles.id`, nil),
	`Roles.Name`:        param.Tf(`roles.name`, nil),
	`Roles.Description`: param.Tf(`roles.description`, nil),
}
