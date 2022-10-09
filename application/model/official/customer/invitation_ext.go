package customer

import "github.com/admpub/webx/application/dbschema"

type InvitationCustomerExt struct {
	*dbschema.OfficialCustomerInvitationUsed
	Customer   *dbschema.OfficialCustomer           `db:"-,relation=id:customer_id|gtZero"`
	Level      *dbschema.OfficialCustomerLevel      `db:"-,relation=id:level_id|gtZero"`
	AgentLevel *dbschema.OfficialCustomerAgentLevel `db:"-,relation=id:agent_level_id|gtZero"`
	RoleList   []*dbschema.OfficialCustomerRole     `db:"-,relation=id:role_ids|split"`
}

type InvitationCustomerWithCode struct {
	*dbschema.OfficialCustomerInvitationUsed
	Customer   *dbschema.OfficialCustomer           `db:"-,relation=id:customer_id|gtZero"`
	Level      *dbschema.OfficialCustomerLevel      `db:"-,relation=id:level_id|gtZero"`
	AgentLevel *dbschema.OfficialCustomerAgentLevel `db:"-,relation=id:agent_level_id|gtZero"`
	RoleList   []*dbschema.OfficialCustomerRole     `db:"-,relation=id:role_ids|split"`
	Invitation *dbschema.OfficialCustomerInvitation `db:"-,relation=id:invitation_id|gtZero"`
}
