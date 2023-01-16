package articlelogic

import (
	"strconv"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func GetCategories(c echo.Context, limit int, categoryType ...string) ([]*dbschema.OfficialCommonCategory, error) {
	return GetCategoriesWithLevel(c, 0, limit, categoryType...)
}

func GetCategoriesWithLevel(c echo.Context, level int, limit int, categoryType ...string) ([]*dbschema.OfficialCommonCategory, error) {
	var ctype string
	if len(categoryType) > 0 {
		ctype = categoryType[0]
	}
	if len(ctype) == 0 {
		ctype = modelArticle.GroupName
	}
	cacheKey := `article.getCategories.` + ctype
	cates, ok := c.Internal().Get(cacheKey).([]*dbschema.OfficialCommonCategory)
	if ok {
		return cates, nil
	}
	cond := db.NewCompounds()
	cond.AddKV(`type`, ctype)
	cond.AddKV(`level`, level)
	cond.AddKV(`disabled`, `N`)
	cateM := official.NewCategory(c)
	_, err := cateM.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, limit, cond.And())
	if err != nil {
		return nil, err
	}
	cates = cateM.Objects()
	c.Internal().Set(cacheKey, cates)
	return cates, nil
}

func GetSubCategories(c echo.Context, parentId int, limit int, categoryType ...string) ([]*dbschema.OfficialCommonCategory, error) {
	var ctype string
	if len(categoryType) > 0 {
		ctype = categoryType[0]
	}
	if len(ctype) == 0 {
		ctype = modelArticle.GroupName
	}
	cacheKey := `article.getSubCategories.` + ctype + `.` + strconv.Itoa(parentId)
	cates, ok := c.Internal().Get(cacheKey).([]*dbschema.OfficialCommonCategory)
	if ok {
		return cates, nil
	}
	cond := db.NewCompounds()
	cond.AddKV(`type`, ctype)
	cond.AddKV(`parent_id`, parentId)
	cond.AddKV(`disabled`, `N`)
	cateM := official.NewCategory(c)
	_, err := cateM.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, limit, cond.And())
	if err != nil {
		return nil, err
	}
	cates = cateM.Objects()
	c.Internal().Set(cacheKey, cates)
	return cates, nil
}
