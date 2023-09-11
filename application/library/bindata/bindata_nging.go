//go:build bindata
// +build bindata

package bindata

import (
	"path"
	"strings"

	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	mwBindata "github.com/webx-top/echo/middleware/bindata"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/manager"

	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/initialize/backend"
	"github.com/admpub/nging/v5/application/library/bindata"
	"github.com/admpub/nging/v5/application/library/modal"
	selfBackend "github.com/admpub/webx/application/initialize/backend"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/xtemplate"
)

func Initialize(callbacks ...func()) {
	bindata.BackendTmplAssetPrefix = "template"
	bindata.FrontendTmplAssetPrefix = "template"
	if len(callbacks) > 0 && callbacks[0] != nil {
		callbacks[0]()
	}
	bindata.Initialize()

	backend.AssetsDir = backend.DefaultAssetsDir
	backend.TemplateDir = backend.DefaultTemplateDir
	backend.RendererDo = func(renderer driver.Driver) {
		selfBackend.TmplPathFixers.SetTmplDir(renderer.TmplDir()).SetHandler(func(c echo.Context, theme string, tmpl string) string {
			var found bool
			tmpl, found = frontend.TmplPathFixers.Fix(c, bindata.BackendTmplAssetFS, theme, tmpl)
			if found {
				return tmpl
			}
			return path.Join(`backend`, tmpl)
		})
		renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
			var theme string
			return selfBackend.TmplPathFixers.Handle(c, theme, tmpl)
		})
	}
	frontend.RendererDo = func(renderer driver.Driver) {
		frontend.TmplPathFixers.SetTmplDir(renderer.TmplDir()).SetHandler(func(c echo.Context, theme string, tmpl string) string {
			var found bool
			tmpl, found = frontend.TmplPathFixers.Fix(c, bindata.FrontendTmplAssetFS, theme, tmpl)
			if found {
				return tmpl
			}
			return path.Join(`frontend`, tmpl)
		})
		renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
			theme := c.Internal().String(`theme`, `default`)
			return frontend.TmplPathFixers.Handle(c, theme, tmpl)
		})
	}

	if echo.String(`LABEL`) != `dev` { // 在开发环境下不启用，避免无法测试 bindata 真实效果
		// 在 bindata 模式，支持优先读取本地的静态资源文件和模板文件，在没有找到的情况下才读取 bindata 内的文件

		// StaticMW

		fileSystems := xtemplate.NewFileSystems()
		fileSystems.Register(xtemplate.NewStaticDir(backend.AssetsDir, "/public/assets/")) // 注册本地文件系统内的文件
		fileSystems.Register(bindata.StaticAssetFS)                                        // 注册 bindata 打包的文件
		bootconfig.StaticMW = mwBindata.Static("/public/assets/", fileSystems)

		// Template file manager

		// 后台
		backendManagers := []driver.Manager{
			manager.New(),             // 本地文件系统内的模板文件
			bootconfig.BackendTmplMgr, // bindata 打包的模板文件
		}
		backendMultiManager := xtemplate.NewMultiManager(backend.TemplateDir, backendManagers...)
		log.Debugf(`%s Enabled MultiManager (num: %d)`, color.GreenString(`[backend.renderer]`), len(backendManagers))
		bootconfig.BackendTmplMgr = backendMultiManager

		// 前台
		frontendManagers := []driver.Manager{
			manager.New(),              // 本地文件系统内的模板文件
			bootconfig.FrontendTmplMgr, // bindata 打包的模板文件
		}
		frontendMultiManager := xtemplate.NewMultiManager(frontend.TemplateDir, frontendManagers...)
		log.Debugf(`%s Enabled MultiManager (num: %d)`, color.GreenString(`[frontend.renderer]`), len(frontendManagers))
		bootconfig.FrontendTmplMgr = frontendMultiManager
	}

	frontend.StaticMW = nil

	modal.PathFixer = func(c echo.Context, file string) string {
		file = strings.TrimPrefix(file, backend.TemplateDir+`/`)
		return selfBackend.TmplPathFixers.Handle(c, ``, file)
	}
}
