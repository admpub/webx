package upload

import (
	"github.com/admpub/nging/v5/application/registry/upload"
)

func init() {
	upload.Subdir.Add(`category`, `分类`)
	upload.Subdir.Add(`navigate`, `导航`)
	upload.Subdir.Add(`friendlink`, `友情链接`)
}
