package index

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/handler/frontend/index/custom"
)

func Custom(c echo.Context) error {
	page := c.Param(`page`)
	h, ok := custom.Pages[page]
	if !ok {
		return echo.ErrNotFound
	}
	err := h.Handle(c)
	tmpl := c.Internal().String(`custom-page-tmpl`)
	if len(tmpl) == 0 {
		tmpl = `index/custom/page-` + page
	}
	return c.Render(tmpl, common.Err(c, err))
}
