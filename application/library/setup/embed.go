package setup

import (
	_ "embed"

	"github.com/admpub/nging/v4/application/library/config"
)

//go:embed install.sql
var installSQL string

func init() {
	config.RegisterInstallSQL(`webx`, installSQL)
}
