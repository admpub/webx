package article

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/article/user"
	"github.com/coscms/webfront/initialize/frontend"
)

func init() {
	frontend.Register(func(g echo.RouteRegister) {
		g.Route(`GET`, `/article/<id:\d+>`, Detail).SetName(`article.detail`)
		g.Route(`GET`, `/article/<id:\d+>/:op`, Detail).SetName(`article.detailWithOp`)
		g.Route(`GET`, `/articlesBy/:sourceTable/:sourceId`, ArticleListBy).SetName(`article.listBy`)
		g.Route(`GET`, `/articles`, List).SetName(`article.list`)
		g.Route(`POST`, `/article/like`, ArticleLike).SetName(`article.like`)
		g.Route(`POST`, `/article/hate`, ArticleHate).SetName(`article.hate`)
		g.Route(`POST`, `/article/comment_add`, ArticleCommentAdd).SetName(`article.comment.add`)
		g.Route(`GET`, `/article/comment_list`, ArticleCommentList).SetName(`article.comment.list`)
		g.Route(`GET`, `/article/comment_reply_list`, ArticleCommentReplyList).SetName(`article.comment.replyList`)
		g.Route(`POST`, `/article/comment_like`, CommentLike).SetName(`article.comment.like`)
		g.Route(`POST`, `/article/comment_hate`, CommentHate).SetName(`article.comment.hate`)
		g.Route(`POST`, `/article/pay/:id`, Pay).SetName(`article.pay`)
		g.Route(`GET,POST`, `/article/redirect`, Redirect).SetName(`article.redirect`)
		g.Route(`GET`, `/article/tags`, Tags).SetName(`article.tags`)
		g.Route(`GET`, `/article`, Index).SetName(`article.index`)
	})
}
