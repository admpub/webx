package article

import (
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
)

var Contype = official.Contype

type ArticleWithOwner struct {
	*dbschema.OfficialCommonArticle
	User     *User     `db:"-,relation=id:owner_id|gtZero|eq(owner_type:user),columns=id&username&avatar" json:",omitempty"`
	Customer *Customer `db:"-,relation=id:owner_id|gtZero|eq(owner_type:customer),columns=id&name&avatar" json:",omitempty"`
	Category *Category `db:"-,relation=id:category_id|gtZero,columns=id&name" json:",omitempty"`
}

type User struct {
	Id       uint   `db:"id"`
	Username string `db:"username"`
	Avatar   string `db:"avatar"`
}

func (u *User) Name_() string {
	return dbschema.WithPrefix(`nging_user`)
}

type Customer struct {
	Id     uint64 `db:"id"`
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
}

func (c *Customer) Name_() string {
	return dbschema.WithPrefix(`official_customer`)
}

type Category struct {
	Id   uint   `db:"id"`
	Name string `db:"name"`
}

func (c *Category) Name_() string {
	return dbschema.WithPrefix(`official_common_category`)
}

type ArticleAndSourceInfo struct {
	*dbschema.OfficialCommonArticle
	Categories []*dbschema.OfficialCommonCategory
	SourceInfo echo.KV `db:"-"`
}

func (a *ArticleAndSourceInfo) GetCategory1() uint {
	return a.Category1
}

func (a *ArticleAndSourceInfo) GetCategory2() uint {
	return a.Category2
}

func (a *ArticleAndSourceInfo) GetCategory3() uint {
	return a.Category3
}

func (a *ArticleAndSourceInfo) GetCategoryID() uint {
	return a.CategoryId
}

func (a *ArticleAndSourceInfo) AddCategory(g *dbschema.OfficialCommonCategory) {
	a.Categories = append(a.Categories, g)
}

func WithSourceInfo(ctx echo.Context, list []*ArticleAndSourceInfo) error {
	sourceTableIds := map[string][]string{}
	sourceIDIndexes := map[string][]int{}
	for i, v := range list {
		if len(v.SourceTable) == 0 || len(v.SourceId) == 0 {
			continue
		}
		if _, ok := sourceTableIds[v.SourceTable]; !ok {
			sourceTableIds[v.SourceTable] = []string{}
		}
		sourceTableIds[v.SourceTable] = append(sourceTableIds[v.SourceTable], v.SourceId)
		scKey := v.SourceTable + `:` + v.SourceId
		if _, ok := sourceIDIndexes[scKey]; !ok {
			sourceIDIndexes[scKey] = []int{}
		}
		sourceIDIndexes[scKey] = append(sourceIDIndexes[scKey], i)
	}
	if len(sourceTableIds) > 0 {
		for sourceTable, sourceIds := range sourceTableIds {
			infoMapGetter := Source.GetInfoMapGetter(sourceTable)
			if infoMapGetter == nil {
				continue
			}
			infoMap, err := infoMapGetter(ctx, sourceIds...)
			if err != nil {
				return err
			}
			if infoMap == nil {
				continue
			}
			for sourceID, info := range infoMap {
				keys, ok := sourceIDIndexes[sourceTable+`:`+sourceID]
				if !ok {
					continue
				}
				for _, index := range keys {
					list[index].SourceInfo = info
				}
			}
		}
	}
	return nil
}
