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
	"github.com/admpub/nging/v4/application/initialize/backend"
	"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/modal"
	"github.com/admpub/nging/v4/application/library/ntemplate"
	selfBackend "github.com/admpub/webx/application/initialize/backend"
	"github.com/admpub/webx/application/initialize/frontend"
)

var (
	StaticOptions = &middleware.StaticOptions{
		Root: "",
		Path: "/public/assets/frontend/",
	}
	NgingDir             = `../nging`
	WebxDir              = `../webx`
	BackendTemplateDirs  = bindata.PathAliases     //{prefix:templateDir}
	FrontendTemplateDirs = ntemplate.PathAliases{} //{prefix:templateDir}
)

func Initialize(callbacks ...func()) {
	backend.AssetsDir = filepath.Join(NgingDir, `public/assets`)
	backend.TemplateDir = filepath.Join(NgingDir, `template/backend`)
	bindata.StaticOptions.AddFallback(filepath.Join(WebxDir, `public/assets`))
	if len(callbacks) > 0 && callbacks[0] != nil {
		callbacks[0]()
	}
	bindata.Initialize()
	backendTemplateDir, err := filepath.Abs(filepath.Join(WebxDir, `template/backend`))
	if err != nil {
		panic(err)
	}
	log.Debug(`[backend] `, `Template subfolder "official" is relocated to: `, backendTemplateDir)
	selfBackend.TmplPathFixers.AddDir(`official`, backendTemplateDir)
	backendUniqueTemplateDirs := map[string]struct{}{}
	for prefix, templateDirs := range BackendTemplateDirs {
		for _, templateDir := range templateDirs {
			log.Debug(`[backend] `, `Template subfolder "`+prefix+`" is relocated to: `, templateDir)
			backendUniqueTemplateDirs[templateDir] = struct{}{}
			selfBackend.TmplPathFixers.AddDir(prefix, templateDir)
		}
	}
	backend.RendererDo = func(renderer driver.Driver) {
		selfBackend.TmplPathFixers.SetTmplDir(renderer.TmplDir()).SetHandler(func(c echo.Context, theme string, tmpl string) string {
			var found bool
			tmpl, found = selfBackend.TmplPathFixers.Fix(c, theme, tmpl)
			if found {
				return tmpl
			}
			return filepath.Join(renderer.TmplDir(), tmpl)
		})
		renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
			var theme string
			return selfBackend.TmplPathFixers.Handle(c, theme, tmpl)
		})
		renderer.Manager().AddWatchDir(backendTemplateDir)
		for templateDir := range backendUniqueTemplateDirs {
			log.Debug(`[backend] `, `Watch folder: `, templateDir)
			renderer.Manager().AddWatchDir(templateDir)
		}
	}
	modal.PathFixer = func(c echo.Context, file string) string {
		fileNew := strings.TrimPrefix(file, backend.TemplateDir+`/`)
		newPath, found := selfBackend.TmplPathFixers.Fix(c, ``, fileNew)
		if found {
			return newPath
		}
		return file
	}
	//注册前台静态资源
	if len(StaticOptions.Root) == 0 {
		StaticOptions.Root = filepath.Join(WebxDir, `public/assets/frontend`)
	}
	frontend.StaticMW = middleware.Static(StaticOptions)
	frontend.TemplateDir = filepath.Join(WebxDir, frontend.DefaultTemplateDir) //模板文件夹
	frontend.AssetsDir = filepath.Join(WebxDir, frontend.DefaultAssetsDir)     //素材文件夹
	frontendTemplateDir := filepath.Join(WebxDir, `template/frontend`)
	FrontendTemplateDirs.AddAllSubdir(frontendTemplateDir)
	//FrontendTemplateDirs.Add(`default`, frontendTemplateDir)
	frontendUniqueTemplateDirs := map[string]struct{}{}
	for prefix, templateDirs := range FrontendTemplateDirs {
		for _, templateDir := range templateDirs {
			frontendUniqueTemplateDirs[templateDir] = struct{}{}
			log.Debug(`[frontend] `, `Template subfolder "`+prefix+`" is relocated to: `, templateDir)
			frontend.TmplPathFixers.AddDir(prefix, templateDir)
		}
	}
	frontend.RendererDo = func(renderer driver.Driver) {
		frontend.TmplPathFixers.SetTmplDir(renderer.TmplDir()).SetHandler(func(c echo.Context, theme string, tmpl string) string {
			var found bool
			tmpl, found = frontend.TmplPathFixers.Fix(c, theme, tmpl)
			if found {
				return tmpl
			}
			if len(theme) > 0 {
				tmpl = theme + `/` + tmpl
			}
			return filepath.Join(renderer.TmplDir(), tmpl)
		})
		renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
			theme := c.Internal().String(`theme`, `default`)
			return frontend.TmplPathFixers.Handle(c, theme, tmpl)
		})
		for _, templateDir := range FrontendTemplateDirs.TmplDirs() {
			log.Debug(`[frontend] `, `Watch folder: `, templateDir)
			renderer.Manager().AddWatchDir(templateDir)
		}
	}
}
