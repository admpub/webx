package article

import (
	"fmt"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/top"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

type (
	ContentHideInfo struct {
		isHide ContentHideDetectorFunc
	}
	cachedContentHideInfo struct {
		isHide    bool
		msgOnHide string
	}
	ContentHideDetectorFunc func(
		customer *dbschema.OfficialCustomer,
		article *dbschema.OfficialCommonArticle,
		hideContent string,
		args ...string,
	) bool
)

func NewContentHideInfo(hideDetector ContentHideDetectorFunc) *ContentHideInfo {
	return &ContentHideInfo{isHide: hideDetector}
}

var ContentHideDetector = echo.NewKVData()

func init() {
	ContentHideDetectorRegister(`signIn`, `登录后才能查看`, func(
		customer *dbschema.OfficialCustomer,
		article *dbschema.OfficialCommonArticle,
		hideContent string,
		args ...string,
	) bool {
		return customer == nil // 返回true表示需要隐藏内容，否则显示内容
	}, `此处内容需要登录后方可阅读`)
	ContentHideDetectorRegister(`owner`, `作者才能查看`, func(
		customer *dbschema.OfficialCustomer,
		article *dbschema.OfficialCommonArticle,
		hideContent string,
		args ...string,
	) bool {
		if customer == nil {
			return true
		}
		if article.OwnerType == `customer` {
			return customer.Id != article.OwnerId
		}
		if customer.Uid < 1 {
			return true
		}
		return uint64(customer.Uid) != article.OwnerId
	}, `此处内容只有本文作者才能查看`)
	ContentHideDetectorRegister(`paid`, `需要付款后才能查看(需设置售价)`, func(
		customer *dbschema.OfficialCustomer,
		article *dbschema.OfficialCommonArticle,
		hideContent string,
		args ...string,
	) bool {
		if article.Price <= 0 {
			return false
		}
		if customer == nil {
			return true
		}
		walletM := modelCustomer.NewWallet(article.Context())
		exists, _ := walletM.Flow.Exists(nil, db.And(
			db.Cond{`customer_id`: customer.Id},
			db.Cond{`asset_type`: `money`},
			db.Cond{`amount_type`: `balance`},
			db.Cond{`source_type`: `buy`},
			db.Cond{`source_table`: `official_common_article`},
			db.Cond{`source_id`: article.Id},
		))
		return !exists
	}, `此处内容需要付款后才能查看`)
}

func ContentHideDetectorRegister(k string, description string, hideDetector ContentHideDetectorFunc, v string) {
	item := echo.NewKV(k, v).SetX(NewContentHideInfo(hideDetector))
	item.SetHKV(`description`, description)
	ContentHideDetector.AddItem(item)
}

func GetContentHideDetector(
	customer *dbschema.OfficialCustomer,
	article *dbschema.OfficialCommonArticle,
) top.HideDetector {
	cached := map[string]*cachedContentHideInfo{}
	return top.HideDetector(func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		if info, ok := cached[hideType]; ok {
			return info.isHide, info.msgOnHide
		}
		defer func() {
			if hideType == `paid` {
				priceStr := fmt.Sprintf(`%0.2f`, article.Price)
				msgOnHide += `<button type="button" class="btn btn-info btn-block tx-white" data-article-payment="` + fmt.Sprint(article.Id) + `" data-article-price="` + priceStr + `">`
				msgOnHide += `<i class="fa fa-money"></i> `
				msgOnHide += article.Context().T(`支付`)
				msgOnHide += ` (<strong>` + priceStr + ` ` + article.Context().T(`元`) + `</strong>)`
				msgOnHide += `</button>`
			}
			cached[hideType] = &cachedContentHideInfo{isHide: hide, msgOnHide: msgOnHide}
		}()
		item := ContentHideDetector.GetItem(hideType)
		if item == nil {
			return
		}
		info, ok := item.X.(*ContentHideInfo)
		if !ok {
			return
		}
		hide = info.isHide(customer, article, hideContent, args...)
		msgOnHide = item.V
		return
	})
}
