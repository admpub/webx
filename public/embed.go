//go:build __misc__

package public

import (
	"embed"

	_ "github.com/admpub/nging/v5/public"
)

//go:embed assets
var assets embed.FS
