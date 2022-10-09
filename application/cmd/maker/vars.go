package maker

import (
	"path/filepath"

	"github.com/webx-top/echo"
)

var (
	// SampleDir 样板代码文件夹
	SampleDir = filepath.Join(echo.Wd(), `application`, `cmd`, `maker`, `sample`)
	//HandlerDir handler保存文件夹
	HandlerDir = filepath.Join(echo.Wd(), `application`, `handler`, `backend`)
	//ModelDir model保存文件夹
	ModelDir = filepath.Join(echo.Wd(), `application`, `model`)
	//TemplateDir 模板保存文件夹
	TemplateDir = filepath.Join(echo.Wd(), `template`, `backend`)
	//DefaultCLIConfig 命令行参数
	DefaultCLIConfig = &CLIConfig{}
)

type CLIConfig struct {
	Tables            string
	Group             string
	DBKey             string
	SwitchableFields  string
	MustHasPrimaryKey bool
}
