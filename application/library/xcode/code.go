package xcode

import (
	"github.com/webx-top/echo/code"
)

var (
	// OrderAlreadyPaid 订单已付款
	OrderAlreadyPaid code.Code = -1008
	// WaitForTheLastSubmissionToComplete 等待上次的提交完成
	WaitForTheLastSubmissionToComplete code.Code = -1009
	// DataAlreadyDeleted 数据已经被删除了
	DataAlreadyDeleted  code.Code = -1010
	InvalidSession      code.Code = -1090
	SignatureHasExpired code.Code = -1091 // 签名已经过期
)

func init() {
	code.Register(OrderAlreadyPaid, `OrderAlreadyPaid`)
	code.Register(WaitForTheLastSubmissionToComplete, `WaitForTheLastSubmissionToComplete`)
	code.Register(DataAlreadyDeleted, `DataAlreadyDeleted`)
	code.Register(InvalidSession, `InvalidSession`)
	code.Register(SignatureHasExpired, `SignatureHasExpired`)
}
