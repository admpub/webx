package comment

import (
	"github.com/coscms/webcore/library/perm"
	"github.com/admpub/webx/application/library/xcommon"
	"github.com/admpub/webx/application/library/xrole"
)

const (
	BehaviorName = `comment`
)

func init() {
	xrole.Behaviors.Register(BehaviorName, `评论设置`,
		perm.BehaviorOptFormHelpBlock(`配置评论发布频率。maxPerDay - 表示每天的最大发布数量(<=0代表禁止发布); maxPending - 表示待审核评论上限(<=0代表不限)`),
		perm.BehaviorOptValue(&xcommon.ConfigCustomerAdd{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &xcommon.ConfigCustomerAdd{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
}
