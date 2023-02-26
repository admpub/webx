package comment

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

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
	modelArticle.ContentHideDetectorRegister(`comment`, `评论后才能查看`, func(params *modelArticle.ContentHideParams) bool {
		if params.Customer == nil {
			return true
		}
		if params.Customer.Uid > 0 { // 管理员不用回复
			return false
		}
		cmtM := NewComment(params.Context)
		exists, _ := cmtM.Exists(nil, db.And(
			db.Cond{`target_type`: modelArticle.GroupName},
			db.Cond{`target_id`: params.Article.Id},
			db.Cond{`owner_type`: `customer`},
			db.Cond{`owner_id`: params.Customer.Id},
		))
		return !exists
	}, `此处内容需要评论后方可阅读`)
}
