//go:build __misc__

package template

import (
	"embed"

	_ "github.com/admpub/nging/v5/template"
)

//go:embed backend frontend
var assets embed.FS
