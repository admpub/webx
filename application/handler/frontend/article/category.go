package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
)

func getCategories(c echo.Context, limit int) ([]*dbschema.OfficialCommonCategory, error) {
	cates, ok := c.Internal().Get(`article.getCategories`).([]*dbschema.OfficialCommonCategory)
	if ok {
		return cates, nil
	}
	cond := db.NewCompounds()
	cond.AddKV(`type`, modelArticle.GroupName)
	cond.AddKV(`level`, 0)
	cond.AddKV(`disabled`, `N`)
	cateM := official.NewCategory(c)
	_, err := cateM.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, limit, cond.And())
	if err != nil {
		return nil, err
	}
	cates = cateM.Objects()
	c.Internal().Set(`article.getCategories`, cates)
	return cates, nil
}
