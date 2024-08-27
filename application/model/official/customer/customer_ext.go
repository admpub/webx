package customer

import (
	"github.com/webx-top/echo/param"

	"github.com/admpub/webx/application/dbschema"
	modelLevel "github.com/admpub/webx/application/model/official/level"
)

type CustomerAndGroup struct {
	*dbschema.OfficialCustomer
	Group  *dbschema.OfficialCommonGroup        `db:"-,relation=id:group_id|gtZero"`
	Levels []*modelLevel.RelationExt            `db:"-,relation=customer_id:id|gtZero" json:",omitempty"`
	Agent  *dbschema.OfficialCustomerAgentLevel `db:"-,relation=id:agent_level|gtZero"`
	Roles  []*dbschema.OfficialCustomerRole     `db:"-,relation=id:role_ids|notEmpty|split"`
}

func (d *CustomerAndGroup) AsMap() param.Store {
	m := d.OfficialCustomer.AsMap()
	if d.Group != nil {
		m[`Group`] = d.Group.AsMap()
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
