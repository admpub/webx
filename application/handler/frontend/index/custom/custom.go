package custom

import "github.com/webx-top/echo"

var Pages = map[string]echo.HandlerFunc{}

func Register(ident string, handle echo.HandlerFunc) {
	Pages[ident] = handle
}
