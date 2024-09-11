package article

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/logic/articlelogic"
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/echo"
)

func getTags(ctx echo.Context, group ...string) ([]*dbschema.OfficialCommonTags, error) {
	return articlelogic.GetTags(ctx, group...)
}

func Tags(ctx echo.Context) error {
	group := ctx.Query(`group`)
	tags, err := getTags(ctx, group)
	if err != nil {
		return err
	}
	ctx.Set(`tagList`, tags)
	return ctx.Render(`article/tags`, common.Err(ctx, err))
}
