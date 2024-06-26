package initialize

import (
	_ "github.com/admpub/nging/v5/application/initialize/manager"
	"github.com/admpub/nging/v5/application/registry/navigate"
	_ "github.com/admpub/webx/application/initialize/backend"
	_ "github.com/admpub/webx/application/initialize/frontend"
)

var nav = &navigate.List{}

var Project = navigate.NewProject(`内容管理`, `webx`, `/official/customer/index`, nav)

func init() {
	navigate.ProjectGet(`nging`).Name = `其它功能`
	navigate.ProjectAdd(1, Project)
}
