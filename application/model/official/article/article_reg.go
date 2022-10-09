package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/model/official"
)

const GroupName = `article`

func init() {
	official.AddClickFlowTarget(GroupName, official.ClickFlowTargetFunc(func(c echo.Context, id interface{}) (func(typ string) error, func() uint64, error) {
		articleM := NewArticle(c)
		err := articleM.Get(nil, `id`, id)
		if err != nil {
			if err == db.ErrNoMoreRows {
				err = c.NewError(code.DataNotFound, `文章不存在`)
			}
			return nil, nil, err
		}
		return func(typ string) error {
			field := typ + `s`
			return articleM.UpdateField(nil, field, db.Raw(field+`+1`), db.Cond{`id`: id})
		}, nil, nil
	}))
}
