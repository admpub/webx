package agent

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/webx/application/dbschema"
)

func NewAgentProduct(ctx echo.Context) *AgentProduct {
	a := &AgentProduct{
		OfficialCustomerAgentProduct: dbschema.NewOfficialCustomerAgentProduct(ctx),
	}
	return a
}

type AgentProduct struct {
	*dbschema.OfficialCustomerAgentProduct
}

func (f *AgentProduct) check() error {
	return nil
}

func (f *AgentProduct) SaveSold(agentID uint64, productID string, productTable string, performance float64) error {
	bean, err := f.GetByProduct(agentID, productID, productTable)
	if err != nil {
		return err
	}
	return f.UpdateFields(nil, echo.H{
		`sold`:        db.Raw(`sold+1`),
		`performance`: db.Raw(`performance+` + param.AsString(performance)),
	}, `id`, bean.Id)
}

func (f *AgentProduct) GetByProduct(agentID uint64, productID string, productTable string) (*dbschema.OfficialCustomerAgentProduct, error) {
	bean := dbschema.NewOfficialCustomerAgentProduct(f.Context())
	err := bean.Get(nil, db.And(
		db.Cond{`agent_id`: agentID},
		db.Cond{`product_id`: productID},
		db.Cond{`product_table`: productTable},
	))
	if err != nil {
		return nil, err
	}
	return bean, nil
}

func (f *AgentProduct) IncrSold(id uint64, n uint) error {
	return f.UpdateFields(nil, echo.H{`sold`: db.Raw(`sold+` + param.AsString(n))}, `id`, id)
}

func (f *AgentProduct) DecrSold(id uint64, n uint) error {
	return f.UpdateFields(nil, echo.H{`sold`: db.Raw(`sold-` + param.AsString(n))}, `id`, id)
}

func (f *AgentProduct) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	if err := f.Exists(f.AgentId, f.ProductId, f.ProductTable); err != nil {
		return nil, err
	}
	return f.OfficialCustomerAgentProduct.Insert()
}

func (f *AgentProduct) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	if err := f.ExistsOther(f.AgentId, f.ProductId, f.ProductTable, f.Id); err != nil {
		return err
	}
	return f.OfficialCustomerAgentProduct.Update(mw, args...)
}

func (f *AgentProduct) Exists(agentID uint64, productID string, productTable string) error {
	exists, err := f.OfficialCustomerAgentProduct.Exists(nil, db.And(
		db.Cond{`agent_id`: agentID},
		db.Cond{`product_id`: productID},
		db.Cond{`product_table`: productTable},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `代理商“%d”已经代理过商品“%d”`, agentID, productID)
	}
	return err
}

func (f *AgentProduct) ExistsOther(agentID uint64, productID string, productTable string, id uint64) error {
	exists, err := f.OfficialCustomerAgentProduct.Exists(nil, db.And(
		db.Cond{`agent_id`: agentID},
		db.Cond{`product_id`: productID},
		db.Cond{`product_table`: productTable},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `代理商“%d”已经代理过商品“%d”`, agentID, productID)
	}
	return err
}
