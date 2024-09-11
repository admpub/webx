package index

import (
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/codec"
	"github.com/coscms/webcore/library/config"
	"github.com/admpub/webx/application/handler/frontend/article"
)

var (
	// UnimplementedHandler 默认未实现处理
	UnimplementedHandler = func(c echo.Context) error {
		return c.String(`Unimplemented`)
	}
	// DefaultIndexHandler 默认首页处理
	DefaultIndexHandler = article.List
	// DefaultSearchHandler 默认搜索处理
	DefaultSearchHandler = UnimplementedHandler
)

func Index(c echo.Context) error {
	defaultHandler := config.Setting(`base`).String(`defaultHandler`)
	if len(defaultHandler) > 0 && defaultHandler != c.Path() {
		return c.Dispatch(defaultHandler).Handle(c)
	}
	return DefaultIndexHandler(c)
}

func Search(c echo.Context) error {
	return DefaultSearchHandler(c)
}

func ErrorCode(c echo.Context) error {
	data := c.Data()
	sortBy := c.Form("sortBy")
	if sortBy == `code` {
		data.SetData(code.CodeDict)
	} else {
		result := map[string]code.Code{}
		for n, v := range code.CodeDict {
			result[v.Text] = n
		}
		data.SetData(result)
	}
	return c.JSON(data)
}

func SecureKey(c echo.Context) error {
	ip := c.RealIP()
	ts := time.Now().Unix()
	sm2Pubkey := codec.DefaultPublicKeyHex()
	data := c.Data()
	data.SetData(echo.H{
		`ip`:  ip,
		`ts`:  ts,
		`key`: sm2Pubkey,
	})
	return c.JSON(data)
}
