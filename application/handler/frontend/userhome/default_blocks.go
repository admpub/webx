package userhome

import (
	userArticle "github.com/admpub/webx/application/handler/frontend/article/user"
	"github.com/admpub/webx/application/handler/frontend/user"
	"github.com/admpub/webx/application/handler/frontend/user/profile"
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webfront/middleware"
	"github.com/coscms/webfront/middleware/sessdata"
	registryUserhome "github.com/coscms/webfront/registry/userhome"
	"github.com/webx-top/echo"
)

func init() {

	registryUserhome.BlockAdd(
		0,
		(&dashboard.Block{
			Title:  `我的文章`,
			Ident:  `article`,
			Tmpl:   `article/user/homepage/partial_article`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			m := Homeowner(ctx)
			return userArticle.ListByCustomer(ctx, m.OfficialCustomer)
		}),
		(&dashboard.Block{
			Title:  `个人资料`,
			Ident:  `profile`,
			Tmpl:   `user/homepage/partial_profile`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			return nil
		}),
		(&dashboard.Block{
			Title:  `我关注的`,
			Ident:  `following`,
			Tmpl:   `user/homepage/partial_following`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			customerID := Homeowner(ctx).Id
			return profile.FollowingBy(ctx, customerID)
		}),
		(&dashboard.Block{
			Title:  `关注我的`,
			Ident:  `followers`,
			Tmpl:   `user/homepage/partial_following`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			customerID := Homeowner(ctx).Id
			return profile.FollowersBy(ctx, customerID)
		}),
		(&dashboard.Block{
			Title:  `我的评论`,
			Ident:  `comments`,
			Tmpl:   `user/homepage/partial_comments`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			//TODO: implement
			return nil
		}).SetHiddenDetector(func(ctx echo.Context) bool {
			return true
		}),
		(&dashboard.Block{
			Title:  `我的收藏`,
			Ident:  `favorites`,
			Tmpl:   `user/homepage/partial_favorites`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			//TODO: implement
			return nil
		}).SetHiddenDetector(func(ctx echo.Context) bool {
			return true
		}),
		(&dashboard.Block{
			Title:  `我的喜欢`,
			Ident:  `likes`,
			Tmpl:   `user/homepage/partial_likes`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			//TODO: implement
			return nil
		}).SetHiddenDetector(func(ctx echo.Context) bool {
			return true
		}),
		(&dashboard.Block{
			Title:  `发送私信`,
			Ident:  `msgsend`,
			Tmpl:   `user/homepage/partial_msgsend`,
			Footer: ``,
		}).SetContentGenerator(middleware.SkipCurrentURLPermCheck(middleware.AuthCheck(echo.HandlerFunc(func(ctx echo.Context) error {
			if ctx.IsPost() {
				if err := sessdata.CheckPerm(ctx, `/user/message/send`); err != nil {
					return err
				}
				m := Homeowner(ctx)
				EndHandler(ctx)
				return user.MessageSend(ctx, m.OfficialCustomer)
			}
			return nil
		}))).Handle),
	)

}
