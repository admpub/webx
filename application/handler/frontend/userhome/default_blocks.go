package userhome

import (
	"github.com/coscms/webcore/registry/dashboard"
	userArticle "github.com/admpub/webx/application/handler/frontend/article/user"
	"github.com/admpub/webx/application/handler/frontend/user"
	"github.com/admpub/webx/application/handler/frontend/user/profile"
	"github.com/admpub/webx/application/middleware"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	registryUserhome "github.com/admpub/webx/application/registry/userhome"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
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
			m := ctx.Get(`info`).(*modelCustomer.CustomerAndGroup)
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
			var customerID uint64
			switch m := ctx.Get(`info`).(type) {
			case *modelCustomer.CustomerAndGroup:
				customerID = m.OfficialCustomer.Id
			case param.Store:
				customerID = m.Uint64(`id`)
			default:
				return echo.ErrNotFound
			}
			return profile.FollowingBy(ctx, customerID)
		}),
		(&dashboard.Block{
			Title:  `关注我的`,
			Ident:  `followers`,
			Tmpl:   `user/homepage/partial_following`,
			Footer: ``,
		}).SetContentGenerator(func(ctx echo.Context) error {
			var customerID uint64
			switch m := ctx.Get(`info`).(type) {
			case *modelCustomer.CustomerAndGroup:
				customerID = m.OfficialCustomer.Id
			case param.Store:
				customerID = m.Uint64(`id`)
			default:
				return echo.ErrNotFound
			}
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
				m := ctx.Get(`info`).(*modelCustomer.CustomerAndGroup)
				return user.MessageSend(ctx, m.OfficialCustomer)
			}
			return nil
		}))).Handle),
	)

}
