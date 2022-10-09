package agent

import "github.com/admpub/webx/application/dbschema"

type AgentProfileExt struct {
	*dbschema.OfficialCustomerAgentProfile
	Recv     *dbschema.OfficialCustomerAgentRecv  `json:",omitempty" db:"-,relation=customer_id:customer_id"`
	Customer *dbschema.OfficialCustomer           `json:",omitempty" db:"-,relation=id:customer_id"`
	Level    *dbschema.OfficialCustomerAgentLevel `json:",omitempty" db:"-,relation=id:apply_level|gtZero"`
}
