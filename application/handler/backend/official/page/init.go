package page

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func init() {
	route.RegisterToGroup(`/official`, func(g echo.RouteRegister) {
		g = g.Group(`/page`)

		// 以下功能暂时没有用到，注释掉
		// // 页面
		// g.Route(`GET,POST`, `/index`, Index)
		// g.Route(`GET,POST`, `/add`, Add)
		// g.Route(`GET,POST`, `/edit`, Edit)
		// g.Route(`GET,POST`, `/delete`, Delete)

		// // 区块
		// g.Route(`GET,POST`, `/block_index`, BlockIndex)
		// g.Route(`GET,POST`, `/block_add`, BlockAdd)
		// g.Route(`GET,POST`, `/block_edit`, BlockEdit)
		// g.Route(`GET,POST`, `/block_delete`, BlockDelete)

		// // 布局
		// g.Route(`GET,POST`, `/layout_index`, LayoutIndex)
		// g.Route(`GET,POST`, `/layout_add`, LayoutAdd)
		// g.Route(`GET,POST`, `/layout_edit`, LayoutEdit)
		// g.Route(`GET,POST`, `/layout_delete`, LayoutDelete)

		// 模板管理
		g.Route(`GET,POST`, `/template_index`, TemplateIndex)
		g.Route(`GET,POST`, `/template_edit`, TemplateEdit)
		g.Route(`GET,POST`, `/template_enable`, TemplateEnable)
		g.Route(`GET,POST`, `/template_config`, TemplateConfig)
	})
}
