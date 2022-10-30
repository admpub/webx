package article

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func getTags(ctx echo.Context) ([]*dbschema.OfficialCommonTags, error) {
	tags, ok := ctx.Internal().Get(`article.getTags`).([]*dbschema.OfficialCommonTags)
	if ok {
		return tags, nil
	}
	tagsM := official.NewTags(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`group`, modelArticle.GroupName)
	_, err := common.NewLister(tagsM, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-num`)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return nil, err
	}
	tags = tagsM.Objects()
	ctx.Internal().Set(`article.getTags`, tags)
	return tags, nil
}

func Tags(ctx echo.Context) error {
	tags, err := getTags(ctx)
	if err != nil {
		return err
	}
	ctx.Set(`tagList`, tags)
	return ctx.Render(`article/tags`, handler.Err(ctx, err))
}
