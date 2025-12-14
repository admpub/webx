package article

import (
	"strings"

	"github.com/coscms/webcore/library/common"
	modelArticle "github.com/coscms/webfront/model/official/article"
	modelComment "github.com/coscms/webfront/model/official/comment"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Index(c echo.Context) error {
	articleM := modelArticle.NewArticle(c)
	var err error
	c.SetFunc(`relationList`, articleM.RelationList)
	c.SetFunc(`queryList`, articleM.QueryList)
	c.SetFunc(`articleListByIds`, func(ids string) []*modelArticle.ArticleWithOwner {
		cond := db.NewCompounds()
		_ids := strings.Split(ids, `,`)
		cond.AddKV(`id`, db.In(_ids))
		articles := articleM.CommonQueryList(cond, len(_ids), 0)
		modelArticle.MultilingualArticlesWithOwner(c, articles)
		return articles
	})
	c.SetFunc(`hotCommentArticles`, func(query string, limit int, offset int) []*modelArticle.ArticleWithOwner {
		rows, _ := hotCommentArticles(c, query, limit, offset)
		return rows
	})
	c.SetFunc(`getArticles`, func() []*modelArticle.ArticleWithOwner {
		c.Request().Form().Set(`pageSize`, `20`)
		cond := db.NewCompounds()
		articles, _ := articleM.ListPageSimple(cond)
		modelArticle.MultilingualArticlesWithOwner(c, articles)
		return articles
	})
	c.Set(`listURL`, c.URLFor(`/articles`))
	return c.Render(`article/index`, common.Err(c, err))
}

func hotCommentArticles(ctx echo.Context, query string, limit int, offset int) ([]*modelArticle.ArticleWithOwner, error) {
	commentM := modelComment.NewComment(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`target_type`, `article`)
	if len(query) > 0 {
		r := common.NewSortedURLValues(query)
		r.ApplyCond(cond)
	}
	targetIDs, err := commentM.GetTargetIDs(cond, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(targetIDs) == 0 {
		return nil, err
	}
	articleM := modelArticle.NewArticle(ctx)
	cond.Reset()
	cond.AddKV(`id`, db.In(targetIDs))
	articles := articleM.CommonQueryList(cond, limit, 0, `-comments`)
	modelArticle.MultilingualArticlesWithOwner(ctx, articles)
	return articles, nil
}
