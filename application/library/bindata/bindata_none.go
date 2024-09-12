//go:build !bindata
// +build !bindata

package bindata

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"

	"github.com/admpub/log"
	selfBackend "github.com/admpub/webx/application/initialize/backend"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/modal"
)

var (
	// StaticOptions 前台static中间件选项
	StaticOptions = &middleware.StaticOptions{
		Root:   "",
		Path:   "",
		MaxAge: bootconfig.HTTPCacheMaxAge,
	}
	NgingDir = `../nging`
	WebxDir  = `../webx`
)

// Initialize 后台和前台模板等素材初始化配置
func Initialize(callbacks ...func()) {
	backend.AssetsDir = filepath.Join(NgingDir, `public/assets`)
	backend.TemplateDir = filepath.Join(NgingDir, `template/backend`)
	bindata.StaticOptions.AddFallback(filepath.Join(WebxDir, `public/assets`))
	if len(callbacks) > 0 && callbacks[0] != nil {
		callbacks[0]()
	}
	bindata.Initialize() // 后台素材初始化配置
	backendTemplateDir, err := filepath.Abs(filepath.Join(WebxDir, `template/backend`))
	if err != nil {
		panic(err)
	}
	log.Debug(`[backend] `, `Template subfolder "official" is relocated to: `, backendTemplateDir)
	selfBackend.TmplPathFixers.AddDir(`official`, backendTemplateDir)

	// 应用后台模块的文件别名分组到后台模板路径修正器
	selfBackend.TmplPathFixers.ApplyAliases()

	backend.RendererDo = func(renderer driver.Driver) {
		selfBackend.TmplPathFixers.Register(renderer, backendTemplateDir)
	}
	modal.PathFixer = func(c echo.Context, file string) string {
		file = strings.TrimPrefix(file, backend.TemplateDir+`/`)
		return selfBackend.TmplPathFixers.Handle(c, ``, file)
	}
	//注册前台静态资源
	if len(StaticOptions.Root) == 0 {
		StaticOptions.Root = filepath.Join(WebxDir, `public/assets/frontend`)
	}
	if len(StaticOptions.Path) == 0 {
		StaticOptions.Path = frontend.Prefix + "/public/assets/frontend/"
	}
	StaticOptions.TrimPrefix = frontend.Prefix
	frontend.StaticMW = middleware.Static(StaticOptions)
	frontend.TemplateDir = filepath.Join(WebxDir, frontend.DefaultTemplateDir) //模板文件夹
	frontend.AssetsDir = filepath.Join(WebxDir, frontend.DefaultAssetsDir)     //素材文件夹
	frontendTemplateDir := filepath.Join(WebxDir, `template/frontend`)
	frontend.TmplPathFixers.PathAliases.AddAllSubdir(frontendTemplateDir)
	//frontend.TmplPathFixers.PathAliases.Add(`default`, frontendTemplateDir)

	// 应用前台模块的文件别名分组到前台模板路径修正器
	frontend.TmplPathFixers.ApplyAliases()

	frontend.RendererDo = func(renderer driver.Driver) {
		frontend.TmplPathFixers.SetEnableTheme(true).Register(renderer)
	}
}
