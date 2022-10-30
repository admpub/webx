package middleware

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/echo"
)

func DefaultRenderDataWrapper(ctx echo.Context, data interface{}) interface{} {
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
