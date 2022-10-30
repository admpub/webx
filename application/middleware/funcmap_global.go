package middleware

import (
	"sync"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/webx/application/library/top"
	"github.com/admpub/webx/application/middleware/sessdata"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var (
	tplFuncMap map[string]interface{}
	tplOnce    sync.Once
)

func initTplFuncMap() {
	tplFuncMap = tplfunc.New()
}

func TplFuncMap() map[string]interface{} {
	tplOnce.Do(initTplFuncMap)
	return tplFuncMap
}

func init() {
	tplfunc.TplFuncMap[`ImageProxyURL`] = sessdata.ImageProxyURL
	tplfunc.TplFuncMap[`ResizeImageURL`] = sessdata.ResizeImageURL
	tplfunc.TplFuncMap[`AbsoluteURL`] = sessdata.AbsoluteURL
	tplfunc.TplFuncMap[`PictureHTML`] = sessdata.PictureWithDefaultHTML
	tplfunc.TplFuncMap[`OutputContent`] = sessdata.OutputContent
	tplfunc.TplFuncMap[`Config`] = config.FromDB
	tplfunc.TplFuncMap[`StarsSlice`] = top.StarsSlice
	tplfunc.TplFuncMap[`StarsSlicex`] = top.StarsSlicex
}
