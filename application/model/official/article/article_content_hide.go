package article

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/top"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	modelLevel "github.com/admpub/webx/application/model/official/level"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type (
	ContentHideInfo struct {
		isHide ContentHideDetectorFunc
	}
	ContentHideParams struct {
		echo.Context
		Customer    *dbschema.OfficialCustomer
		Article     *dbschema.OfficialCommonArticle
		HideContent string
		Args        []string
	}
	cachedContentHideInfo struct {
		isHide    bool
		msgOnHide string
	}
	ContentHideDetectorFunc func(*ContentHideParams) bool
	MessageFuncOnHide       func(*ContentHideParams) string
)

func (c *ContentHideParams) Reset() {
	c.Context = nil
	c.Customer = nil
	c.Article = nil
	c.HideContent = ``
	c.Args = nil
}

func NewContentHideInfo(hideDetector ContentHideDetectorFunc) *ContentHideInfo {
	return &ContentHideInfo{isHide: hideDetector}
}

var poolContentHideParams = sync.Pool{
	New: func() interface{} {
		return &ContentHideParams{}
	},
}
var ContentHideDetector = echo.NewKVData()

func init() {
	ContentHideDetectorRegister(`signIn`, `登录后才能查看`, func(params *ContentHideParams) bool {
		return params.Customer == nil // 返回true表示需要隐藏内容，否则显示内容
	}, `此处内容需要登录后方可阅读`)
	ContentHideDetectorRegister(`level`, `只有特定等级的用户才能查看(多个等级ID用半角逗号“,”分隔)`, func(params *ContentHideParams) bool {
		if params.Customer == nil {
			return true // 返回true表示需要隐藏内容，否则显示内容
		}
		if len(params.Args) == 0 || len(params.Args[0]) == 0 { // 没有指定等级时显示内容
			return false
		}
		levelM := modelCustomer.NewLevel(params.Article.Context())
		levelIDStrings := param.StringSlice(strings.Split(params.Args[0], `,`)).Filter()
		levelIDs := make([]interface{}, len(levelIDStrings))
		for index, levelID := range levelIDStrings {
			levelIDs[index] = levelID
		}
		has, _ := levelM.HasLevel(params.Customer.Id, levelIDs...)
		return !has
	}, `此处内容仅供会员查看`, messageFuncOnHideForLevel)
	ContentHideDetectorRegister(`group`, `只有特定用户组才能查看(多个组ID用半角逗号“,”分隔)`, func(params *ContentHideParams) bool {
		if params.Customer == nil {
			return true // 返回true表示需要隐藏内容，否则显示内容
		}
		if params.Customer.GroupId < 1 {
			return true
		}
		if len(params.Args) == 0 || len(params.Args[0]) == 0 { // 没有指定等级时显示内容
			return false
		}
		groupIDStrings := param.StringSlice(strings.Split(params.Args[0], `,`)).Filter()
		return !com.InSlice(param.AsString(params.Customer.GroupId), groupIDStrings)
	}, ``, messageFuncOnHideForGroup)
	ContentHideDetectorRegister(`role`, `只有特定角色才能查看(多个角色ID用半角逗号“,”分隔)`, func(params *ContentHideParams) bool {
		if params.Customer == nil {
			return true // 返回true表示需要隐藏内容，否则显示内容
		}
		if len(params.Customer.RoleIds) < 1 {
			return true
		}
		if len(params.Args) == 0 || len(params.Args[0]) == 0 { // 没有指定等级时显示内容
			return false
		}
		needRoleIDStrings := param.StringSlice(strings.Split(params.Args[0], `,`)).Filter()
		myRoleIDStrings := param.StringSlice(strings.Split(params.Customer.RoleIds, `,`)).Filter()
		for _, myRoleID := range myRoleIDStrings {
			if com.InSlice(myRoleID, needRoleIDStrings) {
				return false
			}
		}
		return true
	}, ``, messageFuncOnHideForRole)
	ContentHideDetectorRegister(`owner`, `作者才能查看`, func(params *ContentHideParams) bool {
		if params.Customer == nil {
			return true
		}
		if params.Article.OwnerType == `customer` {
			return params.Customer.Id != params.Article.OwnerId
		}
		if params.Customer.Uid < 1 {
			return true
		}
		return uint64(params.Customer.Uid) != params.Article.OwnerId
	}, `此处内容只有本文作者才能查看`)
	ContentHideDetectorRegister(`paid`, `需要付款后才能查看(需设置售价)`, func(params *ContentHideParams) bool {
		if params.Article.Price <= 0 {
			return false
		}
		if params.Customer == nil {
			return true
		}
		walletM := modelCustomer.NewWallet(params.Article.Context())
		exists, _ := walletM.Flow.Exists(nil, db.And(
			db.Cond{`customer_id`: params.Customer.Id},
			db.Cond{`asset_type`: `money`},
			db.Cond{`amount_type`: `balance`},
			db.Cond{`source_type`: `buy`},
			db.Cond{`source_table`: `official_common_article`},
			db.Cond{`source_id`: params.Article.Id},
		))
		return !exists
	}, ``, messageFuncOnHideForPaid)
}

type ContentHideMsgFnKey string

var ctxKeyContentHideParams ContentHideMsgFnKey

func ContentHideDetectorRegister(k string, description string, hideDetector ContentHideDetectorFunc, v string, fn ...MessageFuncOnHide) {
	item := echo.NewKV(k, v).SetX(NewContentHideInfo(hideDetector))
	item.SetHKV(`description`, description)
	if len(fn) > 0 && fn[0] != nil {
		item.SetFn(func(c context.Context) interface{} {
			params := c.Value(ctxKeyContentHideParams).(*ContentHideParams)
			msgOnHide := fn[0](params)
			return msgOnHide
		})
	}
	ContentHideDetector.AddItem(item)
}

func messageFuncOnHideForPaid(params *ContentHideParams) string {
	msgOnHide := params.T(`此处内容需要付款后才能查看`)
	priceStr := fmt.Sprintf(`%0.2f`, params.Article.Price)
	msgOnHide += `<button type="button" class="btn btn-info btn-block tx-white" data-article-payment="` + fmt.Sprint(params.Article.Id) + `" data-article-price="` + priceStr + `">`
	msgOnHide += `<i class="fa fa-money"></i> `
	msgOnHide += params.T(`支付`)
	msgOnHide += ` (<strong>` + priceStr + ` ` + params.T(`元`) + `</strong>)`
	msgOnHide += `</button>`
	return msgOnHide
}

func messageFuncOnHideForLevel(params *ContentHideParams) string {
	msgOnHide := params.T(`此处内容仅供会员查看`)
	if len(params.Args) == 0 || len(params.Args[0]) == 0 {
		return msgOnHide
	}
	levelIDStrings := param.StringSlice(strings.Split(params.Args[0], `,`)).Filter()
	m := modelLevel.NewLevel(params.Context)
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`)
	}, 0, -1, `id`, db.In(levelIDStrings))
	if err == nil {
		rows := m.Objects()
		levelNames := make([]string, len(rows))
		for index, level := range rows {
			levelNames[index] = level.Name
		}
		return params.T(`此处内容仅供等级为%s的会员查看`, strings.Join(levelNames, params.T(`、`)))
	}
	return msgOnHide
}

func messageFuncOnHideForGroup(params *ContentHideParams) string {
	msgOnHide := params.T(`此处内容仅供指定用户组查看`)
	if len(params.Args) == 0 || len(params.Args[0]) == 0 {
		return msgOnHide
	}
	groupIDStrings := param.StringSlice(strings.Split(params.Args[0], `,`)).Filter()
	m := dbschema.NewOfficialCommonGroup(params.Context)
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`)
	}, 0, -1, db.And(
		db.Cond{`type`: `customer`},
		db.Cond{`id`: db.In(groupIDStrings)},
	))
	if err == nil {
		rows := m.Objects()
		groupNames := make([]string, len(rows))
		for index, group := range rows {
			groupNames[index] = group.Name
		}
		return params.T(`此处内容仅供用户组为%s的用户查看`, strings.Join(groupNames, params.T(`、`)))
	}
	return msgOnHide
}

func messageFuncOnHideForRole(params *ContentHideParams) string {
	msgOnHide := params.T(`此处内容仅供指定角色用户查看`)
	if len(params.Args) == 0 || len(params.Args[0]) == 0 {
		return msgOnHide
	}
	roleIDStrings := param.StringSlice(strings.Split(params.Args[0], `,`)).Filter()
	m := dbschema.NewOfficialCustomerRole(params.Context)
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`)
	}, 0, -1, db.Cond{`id`: db.In(roleIDStrings)})
	if err == nil {
		rows := m.Objects()
		roleNames := make([]string, len(rows))
		for index, role := range rows {
			roleNames[index] = role.Name
		}
		return params.T(`此处内容仅供%s查看`, strings.Join(roleNames, params.T(`、`)))
	}
	return msgOnHide
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
		params := poolContentHideParams.Get().(*ContentHideParams)
		params.Context = article.Context()
		params.Customer = customer
		params.Article = article
		params.HideContent = hideContent
		params.Args = args

		hide = info.isHide(params)
		msgOnHide = item.V
		if item.Fn() != nil {
			ctx := context.WithValue(context.Background(), ctxKeyContentHideParams, params)
			msgOnHide = item.Fn()(ctx).(string)
		}
		params.Reset()
		poolContentHideParams.Put(params)
		return
	})
}
