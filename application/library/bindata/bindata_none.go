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
	"github.com/admpub/nging/v5/application/initialize/backend"
	"github.com/admpub/nging/v5/application/library/bindata"
	"github.com/admpub/nging/v5/application/library/modal"
	"github.com/admpub/nging/v5/application/library/ntemplate"
	selfBackend "github.com/admpub/webx/application/initialize/backend"
	"github.com/admpub/webx/application/initialize/frontend"

	frontendLib "github.com/admpub/webx/application/library/frontend"
)

var (
	StaticOptions = &middleware.StaticOptions{
		Root: "",
		Path: "/public/assets/frontend/",
	}
	NgingDir             = `../nging`
	WebxDir              = `../webx`
	BackendTemplateDirs  = bindata.PathAliases        //{prefix:templateDir}
	FrontendTemplateDirs = ntemplate.NewPathAliases() //{prefix:templateDir}
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
	BackendTemplateDirs.Range(func(prefix, templateDir string) error {
		log.Debug(`[backend] `, `Template subfolder "`+prefix+`" is relocated to: `, templateDir)
		selfBackend.TmplPathFixers.AddDir(prefix, templateDir)
		return nil
	})
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
		for _, templateDir := range BackendTemplateDirs.TmplDirs() {
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
	frontendLib.StaticMW = middleware.Static(StaticOptions)
	frontendLib.TemplateDir = filepath.Join(WebxDir, frontendLib.DefaultTemplateDir) //模板文件夹
	frontendLib.AssetsDir = filepath.Join(WebxDir, frontendLib.DefaultAssetsDir)     //素材文件夹
	frontendTemplateDir := filepath.Join(WebxDir, `template/frontend`)
	FrontendTemplateDirs.AddAllSubdir(frontendTemplateDir)
	//FrontendTemplateDirs.Add(`default`, frontendTemplateDir)
	FrontendTemplateDirs.Range(func(prefix, templateDir string) error {
		log.Debug(`[frontend] `, `Template subfolder "`+prefix+`" is relocated to: `, templateDir)
		frontend.TmplPathFixers.AddDir(prefix, templateDir)
		return nil
	})
	frontendLib.RendererDo = func(renderer driver.Driver) {
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
