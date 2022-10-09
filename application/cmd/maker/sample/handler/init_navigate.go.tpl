package {{.PkgName}}

import (
	"github.com/admpub/nging/v3/application/registry/navigate"
	"github.com/admpub/webx/application/initialize/backend"
)

var nav = &navigate.List{
	{
		Display:  true,
		Name:     `影片管理`,
		Action:   `{{.Group}}`,
		Icon:     `list`,
		Children: {{.MakeInitNavigation}},
	},
}

func init() {
	backend.Project.NavList.Add(0, *nav...)
}
