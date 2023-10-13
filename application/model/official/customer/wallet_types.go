package customer

import "github.com/webx-top/echo"

var (
	// AssetTypes 资产类型
	AssetTypes = echo.NewKVData()
	// AmountTypes 金额类型
	AmountTypes = echo.NewKVData()
	// 资金流水记录状态
	FlowStatus = echo.NewKVData()
)

// 资产类型
const (
	AssetTypeMoney      = `money`
	AssetTypeIntegral   = `integral`
	AssetTypeCredit     = `credit`
	AssetTypePoint      = `point`
	AssetTypeGold       = `gold`
	AssetTypeSilver     = `silver`
	AssetTypeCopper     = `copper`
	AssetTypeExperience = `experience`
)

// 金额类型
const (
	AmountTypeBalance = `balance`
	AmountTypeFreeze  = `freeze`
)

// 资金流水记录状态
const (
	//状态(pending-待确认;confirmed-已确认;canceled-已取消)
	FlowStatusPending   = `pending`
	FlowStatusConfirmed = `confirmed`
	FlowStatusCanceled  = `canceled`
)

func init() {
	// 注册资产类型
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeMoney,
		V: `人民币`,
		H: echo.H{
			`icon`:   `icon ion-cash`,
			`bg`:     `warning`,
			`symbol`: `&yen;`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeIntegral,
		V: `积分`,
		H: echo.H{
			`icon`: `icon ion-android-cart`,
			`bg`:   `indigo-light`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeCredit,
		V: `信用分`,
		H: echo.H{
			`icon`:              `icon ion-heart`,
			`bg`:                `pink-light`,
			`comment`:           `满分10`,
			`ignoreAccumulated`: true, // 不支持累计历史值
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypePoint,
		V: `点数`,
		H: echo.H{
			`icon`: `iconfont icon-dengji-11`,
			`bg`:   `gray`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeGold,
		V: `金币`,
		H: echo.H{
			`icon`: `iconfont icon-jinbi2 tx-gold`,
			`bg`:   `superlight`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeSilver,
		V: `银币`,
		H: echo.H{
			`icon`: `iconfont icon-tongbi tx-silver`,
			`bg`:   `gray`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeCopper,
		V: `铜钱`,
		H: echo.H{
			`icon`: `iconfont icon-tongbi tx-copper`,
			`bg`:   `superlight`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeExperience,
		V: `经验`,
		H: echo.H{
			`icon`: `iconfont icon-youxiu`,
			`bg`:   `info`,
		},
	})

	// 注册金额类型
	AmountTypes.Add(AmountTypeBalance, `余额`)
	AmountTypes.Add(AmountTypeFreeze, `冻结`)

	// 注册资金流水数据的状态
	FlowStatus.Add(FlowStatusPending, `待确认`)
	FlowStatus.Add(FlowStatusConfirmed, `已确认`)
	FlowStatus.Add(FlowStatusCanceled, `已取消`)
}

func AssetTypeList() []*echo.KV {
	copied := []*echo.KV{}
	for _, assetType := range AssetTypes.Slice() {
		if assetType.H.Bool(`disabled`) {
			continue
		}
		copied = append(copied, assetType)
	}
	return copied
}

// AssetTypeIsIgnoreAccumulated 判断某种资产类型是否不支持累计历史值
func AssetTypeIsIgnoreAccumulated(assetType string) bool {
	item := AssetTypes.GetItem(assetType)
	if item == nil {
		return false
	}
	return item.H.Bool(`ignoreAccumulated`)
}
