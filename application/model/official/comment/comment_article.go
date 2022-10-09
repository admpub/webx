package comment

import (
	"fmt"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official/article"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func commentArticleGetTarget(ctx echo.Context, targetID uint64) (res echo.H, breadcrumb []echo.KV, detailURL string, err error) {
	articleM := article.NewArticle(ctx)
	err = articleM.Get(nil, `id`, targetID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`文章不存在`)
		}
		return
	}
	res = articleM.AsMap()
	id := fmt.Sprint(targetID)
	breadcrumb = []echo.KV{
		{K: `official/article/index`, V: ctx.T(`资讯列表`)},
		{K: `official/article/edit?id=` + id, V: articleM.Title},
	}
	detailURL = `article/` + id + ctx.DefaultExtension()
	return
}

func commentArticleCheck(ctx echo.Context, f *dbschema.OfficialCommonComment) error {
	if f.TargetId < 1 {
		return ctx.E(`资讯ID无效`)
	}
	newsM := dbschema.NewOfficialCommonArticle(ctx)
	err := newsM.Get(nil, `id`, f.TargetId)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return ctx.E(`您要评论的资讯(ID:%d)不存在`, f.TargetId)
	}
	if newsM.CommentAutoDisplay == `Y` {
		f.Display = `Y`
	} else {
		f.Display = `N`
	}
	f.TargetOwnerId = newsM.OwnerId
	f.TargetOwnerType = newsM.OwnerType
	return nil
}

func commentArticleAfterAdd(ctx echo.Context, f *dbschema.OfficialCommonComment) error {
	if f.ReplyCommentId == 0 { //不包含对评论的回复统计，只包含根评论
		newsM := dbschema.NewOfficialCommonArticle(ctx)
		return newsM.UpdateField(nil, `comments`, db.Raw(`comments+1`), db.Cond{`id`: f.TargetId})
	}
	return nil
}

func commentArticleAfterDelete(ctx echo.Context, f *dbschema.OfficialCommonComment) error {
	newsM := dbschema.NewOfficialCommonArticle(ctx)
	return newsM.Use(f.Trans()).UpdateField(nil, `comments`, db.Raw(`comments-1`), db.Cond{`id`: f.TargetId})
}
