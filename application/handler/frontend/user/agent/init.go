package agent

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webfront/initialize/frontend"
	xMW "github.com/coscms/webfront/middleware"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 代理
		agentG := u.Group(`/agent`, func(h echo.Handler) echo.HandlerFunc {
			return func(c echo.Context) error {
				if c.Request().URL().Path() == u.Prefix()+`/agent/apply` {
					return h.Handle(c)
				}
				customer := xMW.Customer(c)
				if customer.AgentLevel < 1 {
					return c.Redirect(xMW.URLFor(`/user/agent/apply`))
				}
				return h.Handle(c)
			}
		})
		agentG.Route(`GET`, `/index`, AgentIndex)
		agentG.Route(`GET,POST`, `/edit`, AgentEdit)
		agentG.Route(`GET,POST`, `/invited`, InvitedList) // invitations
		agentG.Route(`GET,POST`, `/apply`, AgentApply)
	})

}
