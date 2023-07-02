package index

import (
	"strings"
	"time"

	"github.com/admpub/webx/application/library/cache"
	modelAdvert "github.com/admpub/webx/application/model/official/advert"
	"github.com/webx-top/echo"
)

func Advert(ctx echo.Context) error {
	idents := ctx.Paramx(`idents`).String()
	typ := ctx.Form(`type`)
	key := `advert:` + idents
	data := ctx.Data()
	err := cache.XFunc(ctx, key, data, advert(ctx, data, idents), cache.GenOptions(ctx, 300)...)
	if err != nil {
		data.SetError(err)
	}
	if typ == `html` {
		AdvertsHTML(data)
	}
	return ctx.JSON(data)
}

func AdvertsHTML(data echo.Data) {
	h, y := data.GetData().(echo.H)
	if !y {
		return
	}
	adverts, ok := h[`adverts`]
	if !ok {
		return
	}
	switch ads := adverts.(type) {
	case []*modelAdvert.ItemResponse:
		for _, ad := range ads {
			ad.GenHTML()
		}
	case map[string][]*modelAdvert.ItemResponse:
		for _, adList := range ads {
			for _, ad := range adList {
				ad.GenHTML()
			}
		}
	}
}

func advert(ctx echo.Context, data echo.Data, idents string) func() error {
	return func() error {
		m := modelAdvert.NewAdPosition(ctx)
		var (
			identList []string
			identKeys = make(map[string]struct{})
			list      map[string][]*modelAdvert.ItemResponse
			err       error
		)
		for _, ident := range strings.Split(idents, `,`) {
			ident = strings.TrimSpace(ident)
			if len(ident) == 0 {
				continue
			}
			if _, ok := identKeys[ident]; ok {
				continue
			}
			identKeys[ident] = struct{}{}
			identList = append(identList, ident)
		}
		if len(identList) > 0 {
			list, err = m.GetAdvertsByIdent(identList...)
		}
		if err != nil {
			return err
		}
		adData := echo.H{`refreshedAt`: time.Now().Format(`2006-01-02 15:04:05`)}
		if len(identList) == 1 {
			if list == nil {
				adData.Set(`adverts`, []struct{}{})
			} else {
				for _, rows := range list {
					adData.Set(`adverts`, rows)
				}
			}
		} else {
			if list == nil {
				adData.Set(`adverts`, map[string]struct{}{})
			} else {
				adData.Set(`adverts`, list)
			}
		}
		data.SetData(adData)
		return nil
	}
}
