package comment

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
)

func init() {
	official.AddClickFlowTarget(`comment`, official.ClickFlowTargetFunc(func(c echo.Context, id interface{}) (func(typ string) error, func() uint64, error) {
		m := NewComment(c)
		err := m.Get(nil, `id`, id)
		if err != nil {
			if err == db.ErrNoMoreRows {
				err = c.E(`评论不存在`)
			}
			return nil, nil, err
		}
		return func(typ string) error {
			field := typ + `s`
			return m.UpdateField(nil, field, db.Raw(field+`+1`), db.Cond{`id`: id})
		}, nil, nil
	}))
	modelArticle.ContentHideDetectorRegister(`comment`, `评论后才能查看`, func(
		customer *dbschema.OfficialCustomer,
		article *dbschema.OfficialCommonArticle,
		hideContent string,
		args ...string,
	) bool {
		if customer == nil {
			return true
		}
		if customer.Uid > 0 { // 管理员不用回复
			return false
		}
		cmtM := NewComment(article.Context())
		exists, _ := cmtM.Exists(nil, db.And(
			db.Cond{`target_type`: modelArticle.GroupName},
			db.Cond{`target_id`: article.Id},
			db.Cond{`owner_type`: `customer`},
			db.Cond{`owner_id`: customer.Id},
		))
		return !exists
	}, `此处内容需要评论后方可阅读`)
}
