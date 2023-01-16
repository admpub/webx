package articlelogic

import (
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func GetTags(ctx echo.Context, group ...string) ([]*dbschema.OfficialCommonTags, error) {
	var grp string
	if len(group) > 0 {
		grp = group[0]
	}
	if len(grp) == 0 {
		grp = modelArticle.GroupName
	}
	cacheKey := `article.getTags.` + grp
	tags, ok := ctx.Internal().Get(cacheKey).([]*dbschema.OfficialCommonTags)
	if ok {
		return tags, nil
	}
	tagsM := official.NewTags(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`group`, grp)
	cond.AddKV(`display`, `Y`)
	_, err := common.NewLister(tagsM, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-num`)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return nil, err
	}
	tags = tagsM.Objects()
	ctx.Internal().Set(cacheKey, tags)
	return tags, nil
}
