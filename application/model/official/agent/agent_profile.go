package agent

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewAgentProfile(ctx echo.Context) *AgentProfile {
	a := &AgentProfile{
		OfficialCustomerAgentProfile: dbschema.NewOfficialCustomerAgentProfile(ctx),
	}
	return a
}

type AgentProfile struct {
	*dbschema.OfficialCustomerAgentProfile
}

func (f *AgentProfile) check(edit ...bool) error {
	var editMode bool
	if len(edit) > 0 {
		editMode = edit[0]
	}
	if f.CustomerId < 1 {
		return f.Context().NewError(code.InvalidParameter, `请选择客户`).SetZone(`customerId`)
	}
	if f.ApplyLevel <= 0 {
		return f.Context().NewError(code.InvalidParameter, `请选择代理级别`).SetZone(`applyLevel`)
	}
	bean := dbschema.NewOfficialCustomerAgentProfile(f.Context())
	err := bean.Get(nil, `customer_id`, f.CustomerId)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		if editMode { // 编辑旧数据时却没有找到旧数据
			return f.Context().NewError(code.DataNotFound, `代理商数据不存在`).SetZone(`customerId`)
		}
		return nil // 添加模式时，没有找到已存在的数据则正常
	}
	// 找到已存在的数据
	if !editMode { // 添加模式时，找到已存在的数据则不正常
		return f.Context().NewError(code.DataAlreadyExists, `客户“#%v”的代理商信息已经存在`, f.CustomerId).SetZone(`customerId`)
	}
	return err
}

func (f *AgentProfile) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCustomerAgentProfile.Insert()
}

func (f *AgentProfile) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(true); err != nil {
		return err
	}
	return f.OfficialCustomerAgentProfile.Update(mw, args...)
}

func (f *AgentProfile) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*AgentProfileExt, error) {
	var list []*AgentProfileExt
	_, err := handler.NewLister(f, &list, func(r db.Result) db.Result {
		return r.Relation(`Customer`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			return sel.Columns(dbschema.DBI.OmitSelect(dbschema.NewOfficialCustomer(f.Context()), `safe_pwd`, `salt`, `password`)...)
		}).OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	return list, err
}
