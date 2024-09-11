package wallet

import (
	"github.com/coscms/webcore/registry/dashboard"
)

var RechargePage = dashboard.NewPage(`recharge`, map[string][]string{
	`head`: []string{},
	`body`: []string{},
	`foot`: []string{},
})
