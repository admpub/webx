package formbuilder

import (
	"embed"

	"github.com/coscms/forms/common"
)

//go:embed frontend_templates
var templateFS embed.FS

func init() {
	common.FileSystem.Register(templateFS)
	common.SetTmplDir(`bootstrap4`, `frontend_templates`)
	common.SetTmplDir(`frontend`, `frontend_templates`)
}
