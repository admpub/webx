package wallet

import (
	"github.com/admpub/nging/v5/application/registry/dashboard"
)

var RechargePage = dashboard.NewPage(`recharge`, map[string][]string{
	`head`: []string{},
	`body`: []string{},
	`foot`: []string{},
})
