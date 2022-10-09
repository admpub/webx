package customer

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo/param"
)

type CustomerAndGroup struct {
	*dbschema.OfficialCustomer
	Group *dbschema.OfficialCommonGroup        `db:"-,relation=id:group_id|gtZero"`
	Level *dbschema.OfficialCustomerLevel      `db:"-,relation=id:level_id|gtZero"`
	Agent *dbschema.OfficialCustomerAgentLevel `db:"-,relation=id:agent_level|gtZero"`
	Roles []*dbschema.OfficialCustomerRole     `db:"-,relation=id:role_ids|notEmpty|split"`
}

func (d *CustomerAndGroup) AsMap() param.Store {
	m := d.OfficialCustomer.AsMap()
	if d.Group != nil {
		m[`Group`] = d.Group.AsMap()
	}
	if d.Level != nil {
		m[`Level`] = d.Level.AsMap()
	}
	if d.Agent != nil {
		m[`Agent`] = d.Agent.AsMap()
	}
	if len(d.Roles) > 0 {
		roles := make([]param.Store, len(d.Roles))
		for k, v := range d.Roles {
			roles[k] = v.AsMap()
		}
		m[`Roles`] = roles
	}
	return m
}
