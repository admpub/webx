package article

import (
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/logic/articlelogic"
	"github.com/coscms/webfront/library/xkv"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func getTags(ctx echo.Context, group ...string) ([]*dbschema.OfficialCommonTags, error) {
	return articlelogic.GetTags(ctx, group...)
}

func Tags(ctx echo.Context) error {
	group := ctx.Query(`group`)
	v, _ := xkv.GetValue(ctx, `TAGS_PAGE_SIZE`, `200`, `标签列表页数据量`, `设置前台标签列表页中每页显示的标签数量`)
	var pageSize int
	if len(v) > 0 {
		pageSize = param.AsInt(v)
	}
	if pageSize <= 0 {
		pageSize = 200
	}
	common.SetPagingDefaultSize(ctx, pageSize)
	tags, err := getTags(ctx, group)
	if err != nil {
		return err
	}
	i18nm.GetModelsTranslations(ctx, tags)
	ctx.Set(`tagList`, tags)
	ctx.Set(`tagGroup`, group)
	ctx.Set(`tagGroupName`, official.TagGroups.Get(group))
	return ctx.Render(`article/tags`, common.Err(ctx, err))
}
