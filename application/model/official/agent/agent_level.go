package agent

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func NewAgentLevel(ctx echo.Context) *AgentLevel {
	a := &AgentLevel{
		OfficialCustomerAgentLevel: dbschema.NewOfficialCustomerAgentLevel(ctx),
	}
	return a
}

type AgentLevel struct {
	*dbschema.OfficialCustomerAgentLevel
}

func (f *AgentLevel) Add() (pk interface{}, err error) {
	if err := f.Exists(f.Name); err != nil {
		return nil, err
	}
	f.RoleIds = strings.Trim(f.RoleIds, `,`)
	return f.OfficialCustomerAgentLevel.Insert()
}

func (f *AgentLevel) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.ExistsOther(f.Name, f.Id); err != nil {
		return err
	}
	f.RoleIds = strings.Trim(f.RoleIds, `,`)
	return f.OfficialCustomerAgentLevel.Update(mw, args...)
}

func (f *AgentLevel) Exists(name string) error {
	exists, err := f.OfficialCustomerAgentLevel.Exists(nil, db.Cond{`name`: name})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`等级名称“%s”已经使用过了`, name)
	}
	return err
}

func (f *AgentLevel) ExistsOther(name string, id uint) error {
	exists, err := f.OfficialCustomerAgentLevel.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`等级名称“%s”已经使用过了`, name)
	}
	return err
}
