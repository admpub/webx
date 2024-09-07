package advert

import (
	"strings"
	"time"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/cache"
	"github.com/webx-top/echo"
)

type ItemAndPosition struct {
	*dbschema.OfficialAdItem
	Rendered   string
	AdPosition *dbschema.OfficialAdPosition `db:"-,relation=id:position_id|gtZero"`
}

func (i *ItemAndPosition) GetWidth() uint {
	return i.AdPosition.Width
}

func (i *ItemAndPosition) GetHeight() uint {
	return i.AdPosition.Height
}

func (i *ItemAndPosition) GetURL() string {
	return i.Url
}

func (i *ItemAndPosition) GetContent() string {
	return i.Content
}

func (i *ItemAndPosition) GetContype() string {
	return i.Contype
}

type ItemAndRendered struct {
	*dbschema.OfficialAdItem
	Rendered string
}

type CachedAdvert struct {
	List        PositionAdverts
	RefreshedAt time.Time
}

func (c *CachedAdvert) GenHTML() *CachedAdvert {
	c.List.GenHTML()
	return c
}

func GetCachedAdvert(ctx echo.Context, idents ...string) (*CachedAdvert, error) {
	res := &CachedAdvert{}
	if len(idents) > 0 {
		key := `advert:` + strings.Join(idents, `,`)
		err := cache.XFunc(ctx, key, res, func() error {
			m := NewAdPosition(ctx)
			var err error
			res.List, err = m.GetAdvertsByIdent(idents...)
			if err != nil {
				return err
			}
			res.RefreshedAt = time.Now()
			return err
		}, cache.GenOptions(ctx, 300)...)
		if err != nil {
			return nil, err
		}
	}
	if res.List == nil {
		res.List = PositionAdverts{}
	}
	echo.Dump(res)
	return res, nil
}

func GetAdvertForHTML(ctx echo.Context, idents ...string) interface{} {
	sz := len(idents)
	if sz < 1 || (sz == 1 && len(idents[0]) == 0) {
		return ItemsResponse{}
	}
	cc, err := GetCachedAdvert(ctx, idents...)
	if err != nil {
		return ItemsResponse{
			{Rendered: err.Error()},
		}
	}
	if sz == 1 {
		for _, item := range cc.List {
			return item
		}
		return ItemsResponse{}
	}
	return cc.List
}
