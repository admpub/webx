package customer

import (
	"github.com/webx-top/echo"
)

var (
	ComplaintProcesses = echo.NewKVData()
)

func init() {
	ComplaintProcessAdd(`idle`, `空闲`)
	ComplaintProcessAdd(`reject`, `驳回`)
	ComplaintProcessAdd(`queue`, `处理中`)
	ComplaintProcessAdd(`done`, `已处理`)
}

func ComplaintProcessAdd(key, name string) {
	item := &echo.KV{
		K: key,
		V: name,
	}
	ComplaintProcesses.AddItem(item)
}

func ComplaintProcessList() []*echo.KV {
	return ComplaintProcesses.Slice()
}
