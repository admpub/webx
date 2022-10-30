package middleware

import (
	"sync"

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
