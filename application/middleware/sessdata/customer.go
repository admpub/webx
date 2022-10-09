package sessdata

import (
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

// Customer 前台会员客户信息
func Customer(c echo.Context) *dbschema.OfficialCustomer {
	customer, _ := c.Internal().Get(`customer`).(*dbschema.OfficialCustomer)
	return customer
}

// AgentLevel 前台会员客户代理登记信息
func AgentLevel(c echo.Context) *dbschema.OfficialCustomerAgentLevel {
	agentLevel, _ := c.Internal().Get(`agentLevel`).(*dbschema.OfficialCustomerAgentLevel)
	return agentLevel
}

// IsAdmin 是否后台管理员
func IsAdmin(c echo.Context) bool {
	return AdminUID(c) > 0
}

// AdminUID 后台管理员用户ID
func AdminUID(c echo.Context) uint {
	if user := User(c); user != nil {
		return user.Id
	}
	if customer := Customer(c); customer != nil {
		return customer.Uid
	}
	return 0
}
