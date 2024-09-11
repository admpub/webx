package setup

import (
	_ "embed"

	"github.com/coscms/webcore/library/config"
)

//go:embed install.sql
var installSQL string

func init() {
	config.RegisterInstallSQL(`webx`, installSQL)
}
