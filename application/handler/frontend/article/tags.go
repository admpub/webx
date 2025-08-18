package article

import (
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/logic/articlelogic"
	"github.com/coscms/webfront/model/official"
	"github.com/webx-top/echo"
)

func getTags(ctx echo.Context, group ...string) ([]*dbschema.OfficialCommonTags, error) {
	return articlelogic.GetTags(ctx, group...)
}

func Tags(ctx echo.Context) error {
	group := ctx.Query(`group`)
	common.SetPageDefaultSize(ctx, 180)
	tags, err := getTags(ctx, group)
	if err != nil {
		return err
	}
	ctx.Set(`tagList`, tags)
	ctx.Set(`tagGroup`, group)
	ctx.Set(`tagGroupName`, official.TagGroups.Get(group))
	return ctx.Render(`article/tags`, common.Err(ctx, err))
}
