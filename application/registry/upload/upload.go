package upload

import (
	"github.com/coscms/webcore/registry/upload"
)

func init() {
	upload.Subdir.Add(`category`, `分类`)
	upload.Subdir.Add(`navigate`, `导航`)
	upload.Subdir.Add(`friendlink`, `友情链接`)
	upload.Subdir.Add(`membership`, `会员套餐`)
}
