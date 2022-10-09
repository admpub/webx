package agent

import "github.com/webx-top/echo"

const (
	ProfileStatusIdle             = `idle`
	ProfileStatusPending          = `pending`
	ProfileStatusPaid             = `paid`
	ProfileStatusUnconfirm        = `unconfirm`
	ProfileStatusSuccess          = `success`
	ProfileStatusReject           = `reject`
	ProfileStatusCheat            = `cheat`
	ProfileStatusSignedContract   = `signedContract`
	ProfileStatusUnsignedContract = `unsignedContract`
)

// ProfileStatus idle:空闲/草稿;pending:待付款;paid:已付款;unconfirm:未确认;success:申请成功;reject:拒绝;cheat:作弊封号;unsignedContract:未签合同;signedContract:已签合同
var ProfileStatus = echo.NewKVData().
	Add(ProfileStatusIdle, `草稿`).
	Add(ProfileStatusPending, `待付款`,
		echo.KVOptHKV(`applyDescrition`, `如果付款后状态依然没有改变，请联系客服`)).
	Add(ProfileStatusUnconfirm, `未确认`,
		echo.KVOptHKV(`applyDescrition`, `正在等待管理员确认您的资料`)).
	Add(ProfileStatusPaid, `已付款`,
		echo.KVOptHKV(`color`, `primary`),
		echo.KVOptHKV(`applyDescrition`, `已经收到您的付款，请静待管理员确认`)).
	Add(ProfileStatusSuccess, `申请成功`,
		echo.KVOptHKV(`color`, `success`),
		echo.KVOptHKV(`applyDescrition`, `🎉恭喜您成为代理`)).
	Add(ProfileStatusReject, `拒绝`,
		echo.KVOptHKV(`color`, `danger`),
		echo.KVOptHKV(`applyDescrition`, `因资料不符合要求，已经驳回了您的申请，请修改资料后重新提交`)).
	Add(ProfileStatusCheat, `作弊封号`,
		echo.KVOptHKV(`color`, `warning`),
		echo.KVOptHKV(`applyDescrition`, `因存在违规行为，你的代理资格已冻结`)).
	Add(ProfileStatusUnsignedContract, `未签合约`,
		echo.KVOptHKV(`color`, `warning`),
		echo.KVOptHKV(`applyDescrition`, `请签好书面合同后寄回`)).
	Add(ProfileStatusSignedContract, `已签合约`,
		echo.KVOptHKV(`color`, `info`),
		echo.KVOptHKV(`applyDescrition`, `已经签订合同`))
