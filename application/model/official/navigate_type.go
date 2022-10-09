package official

import (
	"github.com/webx-top/echo"
)

var (
	// NavigateTypes 分类的类别
	NavigateTypes = echo.NewKVData()
)

func init() {
	NavigateTypes.Add(`default`,`前台菜单`)
	NavigateTypes.Add(`userCenter`,`用户中心菜单`)
}

// AddNavigateType 登记新的类别
func AddNavigateType(value, label string) {
	NavigateTypes.Add(value, label)
}

// GetAddNavigateTypes 获取类别列表
func GetAddNavigateTypes() []*echo.KV {
	return NavigateTypes.Slice()
}