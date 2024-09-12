package application

import (
	"github.com/admpub/webx/application/handler/backend/official/advert"
	"github.com/admpub/webx/application/handler/backend/official/agent"
	"github.com/admpub/webx/application/handler/backend/official/api"
	"github.com/admpub/webx/application/handler/backend/official/article"
	"github.com/admpub/webx/application/handler/backend/official/customer"
	"github.com/admpub/webx/application/handler/backend/official/page"
	"github.com/admpub/webx/application/handler/backend/official/shorturl"
	"github.com/admpub/webx/application/handler/backend/official/tags"
	"github.com/coscms/webcore/registry/navigate"
	_ "github.com/coscms/webfront/initialize"
)

func Initialize() {
	navigate.ProjectGet(`webx`).NavList.Add(-1,
		article.LeftNavigate,
		tags.LeftNavigate,
		customer.LeftNavigate,
		agent.LeftNavigate,
		api.LeftNavigate,
		shorturl.LeftNavigate,
		page.LeftNavigate,
		advert.LeftNavigate,
	)
}
