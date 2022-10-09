package official

import "github.com/webx-top/echo"

var (
	Contype = echo.NewKVData()
)

func init() {
	Contype.Add(`text`, `纯文本`)
	Contype.Add(`html`, `HTML`)
	Contype.Add(`markdown`, `Markdown`)
}
