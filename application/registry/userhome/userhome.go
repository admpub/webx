package userhome

import (
	"github.com/coscms/webcore/registry/dashboard"
	"github.com/webx-top/echo"
)

var blocks = dashboard.Blocks{}

func BlockRegister(block ...*dashboard.Block) {
	blocks.Add(-1, block...)
}

func BlockAdd(index int, block ...*dashboard.Block) {
	blocks.Add(index, block...)
}

// BlockRemove 删除元素
func BlockRemove(index int) {
	blocks.Remove(index)
}

// BlockSet 设置元素
func BlockSet(index int, list ...*dashboard.Block) {
	blocks.Set(index, list...)
}

func BlockAll(ctx echo.Context) dashboard.Blocks {
	result := make(dashboard.Blocks, len(blocks))
	for k, v := range blocks {
		val := *v
		val.SetHidden(val.IsHidden(ctx))
		result[k] = &val
	}
	return result
}
