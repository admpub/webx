package agent

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func init() {
	//official_customer_agent_profile
	dbschema.DBI.On(`deleting`, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialCustomerAgentProfile)
		// 设置客户代理商等级
		custM := dbschema.NewOfficialCustomer(fm.Context())
		custM.CPAFrom(fm)
		err := custM.UpdateFields(nil, echo.H{`agent_level`: 0}, `id`, fm.CustomerId)
		if err != nil {
			return err
		}

		return err
	}, `official_customer_agent_profile`)
	dbschema.DBI.On(`created,updated`, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialCustomerAgentProfile)
		// 设置客户代理商等级
		custM := dbschema.NewOfficialCustomer(fm.Context())
		custM.CPAFrom(fm)
		var levelID uint
		if fm.Status == `success` {
			levelID = fm.ApplyLevel
		}
		err := custM.UpdateFields(nil, echo.H{`agent_level`: levelID}, `id`, fm.CustomerId)
		if err != nil {
			return err
		}

		return err
	}, `official_customer_agent_profile`)
}
