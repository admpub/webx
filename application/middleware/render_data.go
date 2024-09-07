package middleware

import (
	"html/template"

	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/license"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/logic/articlelogic"
	"github.com/admpub/webx/application/model/official"
	modelAdvert "github.com/admpub/webx/application/model/official/advert"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func init() {
	tplfunc.TplFuncMap[`Advert`] = func(idents ...string) interface{} {
		ctx := defaults.AcquireMockContext()
		r := modelAdvert.GetAdvertForHTML(ctx, idents...)
		defaults.ReleaseMockContext(ctx)
		return r
	}
}

var DefaultRenderDataWrapper = func(ctx echo.Context, data interface{}) interface{} {
	return NewRenderData(ctx, data)
}

func NewRenderData(ctx echo.Context, data interface{}) *RenderData {
	return &RenderData{
		ctx:        ctx,
		RenderData: echo.NewRenderData(ctx, data),
	}
}

type RenderData struct {
	ctx echo.Context
	*echo.RenderData
}

func (r *RenderData) Customer() *dbschema.OfficialCustomer {
	return Customer(r.ctx)
}

func (r *RenderData) CustomerDetail() *modelCustomer.CustomerAndGroup {
	return CustomerDetail(r.ctx)
}

func (r *RenderData) Friendlink(limit int, categoryIds ...uint) []*dbschema.OfficialCommonFriendlink {
	m := official.NewFriendlink(r.ctx)
	list, _ := m.ListShowAndRecord(limit, categoryIds...)
	return list
}

func (r *RenderData) FrontendNav(parentIDs ...uint) []*official.NavigateExt {
	return NavigateList(r.ctx, dbschema.NewOfficialCommonNavigate(r.ctx), `default`, parentIDs...)
}

func (r *RenderData) CustomerNav(parentIDs ...uint) []*official.NavigateExt {
	return NavigateList(r.ctx, dbschema.NewOfficialCommonNavigate(r.ctx), `userCenter`, parentIDs...)
}

func (r *RenderData) SQLQuery() *common.SQLQuery {
	return common.NewSQLQuery(r.ctx)
}

func (r *RenderData) SQLQueryLimit(offset int, limit int, linkID ...int) *common.SQLQuery {
	return common.NewSQLQueryLimit(r.ctx, offset, limit, linkID...)
}

func (r *RenderData) CaptchaForm(tmpl string, args ...interface{}) template.HTML {
	return common.CaptchaForm(r.ctx, tmpl, args...)
}

func (r *RenderData) CaptchaFormWithType(typ string, tmpl string, args ...interface{}) template.HTML {
	return common.CaptchaFormWithType(r.ctx, typ, tmpl, args...)
}

func (r *RenderData) TagList(group ...string) []*dbschema.OfficialCommonTags {
	list, _ := articlelogic.GetTags(r.ctx, group...)
	return list
}

func (r *RenderData) CategoryList(limit int, ctype ...string) []*dbschema.OfficialCommonCategory {
	categories, _ := articlelogic.GetCategories(r.ctx, limit, ctype...)
	return categories
}

func (r *RenderData) SubCategoryList(parentId int, limit int, ctype ...string) []*dbschema.OfficialCommonCategory {
	categories, _ := articlelogic.GetSubCategories(r.ctx, parentId, limit, ctype...)
	return categories
}

func (r *RenderData) SoftwareURL() string {
	if license.SkipLicenseCheck {
		return ``
	}
	return license.ProductURL()
}

func (r *RenderData) SkipLicenseCheck() bool {
	return license.SkipLicenseCheck
}

func (r *RenderData) SoftwareName() string {
	return bootconfig.SoftwareName
}

func (r *RenderData) Advert(idents ...string) interface{} {
	return modelAdvert.GetAdvertForHTML(r.ctx, idents...)
}
