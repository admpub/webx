package frontend

import (
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
)

const (
	Name                  = `frontend`
	DefaultTemplateDir    = `./template/` + Name
	DefaultAssetsDir      = `./public/assets`
	DefaultAssetsURLPath  = `/public/assets/frontend`
	RouteDefaultExtension = `.html`
)

var (
	Prefix             string
	StaticMW           interface{}
	TemplateDir        = DefaultTemplateDir //模板文件夹
	AssetsDir          = DefaultAssetsDir   //素材文件夹
	AssetsURLPath      = DefaultAssetsURLPath
	StaticRootURLPath  = `/public/`
	RendererDo         = func(driver.Driver) {}
	DefaultMiddlewares = []interface{}{middleware.Log()}
)
