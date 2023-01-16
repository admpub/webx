package article

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/logic/articlelogic"
	"github.com/webx-top/echo"
)

func getTags(ctx echo.Context) ([]*dbschema.OfficialCommonTags, error) {
	return articlelogic.GetTags(ctx)
}

func Tags(ctx echo.Context) error {
	tags, err := getTags(ctx)
	if err != nil {
		return err
	}
	ctx.Set(`tagList`, tags)
	return ctx.Render(`article/tags`, handler.Err(ctx, err))
}
