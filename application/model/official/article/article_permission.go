package article

import (
	"github.com/coscms/webcore/library/perm"
	"github.com/admpub/webx/application/library/xcommon"
	"github.com/admpub/webx/application/library/xrole"
)

const (
	BehaviorName = `article`
)

func init() {
	xrole.Behaviors.Register(BehaviorName, `文章投稿设置`,
		perm.BehaviorOptFormHelpBlock(`配置文章投稿频率。maxPerDay - 表示每天的最大发布数量(<=0代表禁止发布); maxPending - 表示待审核文章上限(<=0代表不限)`),
		perm.BehaviorOptValue(&xcommon.ConfigCustomerAdd{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &xcommon.ConfigCustomerAdd{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
}
