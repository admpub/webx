package agent

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/admpub/webx/application/initialize/frontend"
	xMW "github.com/admpub/webx/application/middleware"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		// 代理
		agentG := u.Group(`/agent`, func(h echo.Handler) echo.HandlerFunc {
			return func(c echo.Context) error {
				if c.Request().URL().Path() == handler.FrontendPrefix+`/user/agent/apply` {
					return h.Handle(c)
				}
				customer := xMW.Customer(c)
				if customer.AgentLevel < 1 {
					return c.Redirect(xMW.URLFor(`/user/agent/apply`))
				}
				return h.Handle(c)
			}
		})
		agentG.Route(`GET`, ``, AgentIndex)
		agentG.Route(`GET,POST`, `/edit`, AgentEdit)
		agentG.Route(`GET,POST`, `/invited`, InvitedList) // invitations
		agentG.Route(`GET,POST`, `/apply`, AgentApply)
	})

}
