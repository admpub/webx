package xroleupload

import (
	"sync"

	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/null"
	"github.com/admpub/webx/application/library/xrole"
	sizeBytes "github.com/webx-top/echo/middleware/bytes"
)

type CustomerUpload struct {
	MaxTotalNum       uint64 `json:"maxTotalNum" xml:"maxTotalNum"`   // 最大总数量
	MaxTotalSize      string `json:"maxTotalSize" xml:"maxTotalSize"` // 最大总尺寸
	maxTotalSizeBytes null.Uint64
	mu                sync.RWMutex
}

func (c *CustomerUpload) Combine(source interface{}) interface{} {
	src := source.(*CustomerUpload)
	if src.MaxTotalNum > c.MaxTotalNum {
		c.MaxTotalNum = src.MaxTotalNum
	}
	if ParseSizeBytes(src.MaxTotalSize) > ParseSizeBytes(c.MaxTotalSize) {
		c.MaxTotalSize = src.MaxTotalSize
	}
	return c
}

func (c *CustomerUpload) MaxTotalSizeBytes() uint64 {
	c.mu.RLock()
	maxTotalSizeBytes := c.maxTotalSizeBytes
	c.mu.RUnlock()
	if maxTotalSizeBytes.Valid {
		return maxTotalSizeBytes.Uint64
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxTotalSizeBytes.Uint64 = ParseSizeBytes(c.MaxTotalSize)
	c.maxTotalSizeBytes.Valid = true
	return c.maxTotalSizeBytes.Uint64
}

const BehaviorName = `upload`

func init() {
	xrole.Behaviors.Register(BehaviorName, `上传文件`,
		perm.BehaviorOptFormHelpBlock(`配置上传文件限制。maxTotalNum - 指定客户允许存储的文件最大数量; maxTotalSize - 指定客户允许存储的空间大小(支持的单位有:KB/MB/GB/TB/TB/B),不指定单位时代表单位为B(字节)`),
		perm.BehaviorOptValue(&CustomerUpload{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &CustomerUpload{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
}

func ParseSizeBytes(val string) uint64 {
	v, _ := sizeBytes.Parse(val)
	if v < 0 {
		return 0
	}
	return uint64(v)
}
