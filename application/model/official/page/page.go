package page

import (
	"encoding/json"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func New(ctx echo.Context) *Page {
	m := &Page{
		OfficialPage: dbschema.NewOfficialPage(ctx),
	}
	return m
}

type Page struct {
	*dbschema.OfficialPage
}

func (f *Page) LayoutList() ([]*LayoutWithBlock, error) {
	m := dbschema.NewOfficialPageLayout(f.Context())
	m.Use(f.Trans())
	list := []*LayoutWithBlock{}
	_, err := m.ListByOffset(&list, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, -1, db.And(
		db.Cond{`disabled`: `N`},
		db.Cond{`page_id`: f.Id},
	))
	if err != nil {
		return list, err
	}
	for k, v := range list {
		v.Configs = echo.H{}
		v.ItemConfigs = echo.H{}
		if len(v.OfficialPageLayout.Configs) > 0 {
			err = json.Unmarshal(com.Str2bytes(v.OfficialPageLayout.Configs), &v.Configs)
			if err != nil {
				return list, err
			}
		}
		if v.Block != nil {
			if len(v.Block.ItemConfigs) > 0 {
				err = json.Unmarshal(com.Str2bytes(v.Block.ItemConfigs), &v.ItemConfigs)
				if err != nil {
					return list, err
				}
			}
		}
		list[k] = v
	}
	return list, err
}
