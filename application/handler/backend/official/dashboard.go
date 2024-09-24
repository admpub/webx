package official

import (
	"fmt"
	"html/template"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/dbschema"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func init() {

	httpserver.Backend.Dashboard.Cards.Register(
		(&dashboard.Card{
			IconName:  `fa-users`,
			IconColor: `warning`,
			Short:     `CUSTOMER`,
			Name:      `客户数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			custMdl := modelCustomer.NewCustomer(ctx)
			total, _ := custMdl.Count(nil)
			agents, _ := custMdl.Count(nil, db.Cond{`agent_level`: db.NotEq(0)})
			if agents < 1 {
				return total
			}
			return template.HTML(fmt.Sprintf(`%d <a class="label label-danger" href="%s">%s:%d</a>`, total, `javascript:;`, ctx.T(`代理商`), agents))
		}),
		(&dashboard.Card{
			IconName:  `fa-leaf`,
			IconColor: `success`,
			Short:     `ARTICLES`,
			Name:      `文章数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			m := dbschema.NewOfficialCommonArticle(ctx)
			n, _ := m.Count(nil)
			return n
		}),
		(&dashboard.Card{
			IconName:  `fa-comments`,
			IconColor: `primary`,
			Short:     `COMMENT`,
			Name:      `评论数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			m := dbschema.NewOfficialCommonComment(ctx)
			total, _ := m.Count(nil)
			pendings, _ := m.Count(nil, db.Cond{`display`: `N`})
			if pendings < 1 {
				return total
			}
			return template.HTML(fmt.Sprintf(`%d <a class="label label-danger" href="%s">%s:%d</a>`, total, `/official/article/comment/index?display=N`, ctx.T(`待审`), pendings))
		}),
		(&dashboard.Card{
			IconName:  `fa-warning`,
			IconColor: `danger`,
			Short:     `COMPLAINT`,
			Name:      `客户投诉`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			m := dbschema.NewOfficialCommonComplaint(ctx)
			total, _ := m.Count(nil)
			pendings, _ := m.Count(nil, db.Cond{`process`: `idle`})
			if pendings < 1 {
				return total
			}
			return template.HTML(fmt.Sprintf(`%d <a class="label label-danger" href="%s">%s:%d</a>`, total, `/official/customer/complaint/index?process=idle`, ctx.T(`待处理`), pendings))
		}),
	)
}
