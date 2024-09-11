package shorturl

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/short_url/index`, Index)
		g.Route(`GET,POST`, `/short_url/add`, Add)
		g.Route(`GET,POST`, `/short_url/edit/:id`, Edit)
		g.Route(`GET,POST`, `/short_url/delete/:id`, Delete)
		g.Route(`GET,POST`, `/short_url/analysis`, Analysis)
		g.Route(`GET,POST`, `/short_url/domain_index`, DomainIndex)
		g.Route(`GET,POST`, `/short_url/domain_edit/:id`, DomainEdit)
		g.Route(`GET,POST`, `/short_url/domain_delete/:id`, DomainDelete)
	})
}
