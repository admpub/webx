package agent

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewAgentRecv(ctx echo.Context) *AgentRecv {
	a := &AgentRecv{
		OfficialCustomerAgentRecv: dbschema.NewOfficialCustomerAgentRecv(ctx),
	}
	return a
}

type AgentRecv struct {
	*dbschema.OfficialCustomerAgentRecv
}

func (f *AgentRecv) check(edit ...bool) error {
	var editMode bool
	if len(edit) > 0 {
		editMode = edit[0]
	}
	if f.CustomerId < 1 {
		return f.Context().NewError(code.InvalidParameter, `请选择客户`).SetZone(`customerId`)
	}
	item := RecvMoneyMethod.GetItem(f.RecvMoneyMethod)
	if item == nil {
		return f.Context().NewError(code.InvalidParameter, `暂不支持这种收款方式: %v`, f.RecvMoneyMethod).SetZone(`recvMoneyMethod`)
	}
	if !item.H.Bool(`online`) {
		if len(f.RecvMoneyBranch) == 0 {
			return f.Context().NewError(code.InvalidParameter, `请输入银行支行名称`).SetZone(`recvMoneyBranch`)
		}
	}
	bean := dbschema.NewOfficialCustomerAgentRecv(f.Context())
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

func (f *AgentRecv) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCustomerAgentRecv.Insert()
}

func (f *AgentRecv) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(true); err != nil {
		return err
	}
	return f.OfficialCustomerAgentRecv.Update(mw, args...)
}
