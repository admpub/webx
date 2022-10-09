package wallet

import (
	"github.com/admpub/nging/v4/application/registry/dashboard"
)

var RechargePage = dashboard.NewPage(`recharge`, map[string][]string{
	`head`: []string{},
	`body`: []string{},
	`foot`: []string{},
})
