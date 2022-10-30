package initialize

import (
	"github.com/admpub/nging/v5/application/registry/navigate"
	_ "github.com/admpub/webx/application/initialize/backend"
	_ "github.com/admpub/webx/application/initialize/frontend"
)

var nav = &navigate.List{}

var Project = navigate.NewProject(`Webx`, `webx`, `/official/customer/index`, nav)

func init() {
	navigate.ProjectAdd(1, Project)
}
