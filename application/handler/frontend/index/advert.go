package index

import (
	modelAdvert "github.com/coscms/webfront/model/official/advert"
	"github.com/webx-top/echo"
)

// /advert/ident or /advert/ident1,ident2
func Advert(ctx echo.Context) error {
	idents := ctx.Paramx(`idents`).Split(`,`).Unique().Filter().String()
	typ := ctx.Form(`type`)
	data := ctx.Data()
	cc, err := modelAdvert.GetCachedAdvert(ctx, idents...)
	if err != nil {
		data.SetError(err)
	} else {
		adData := echo.H{`refreshedAt`: cc.RefreshedAt.Format(`2006-01-02 15:04:05`)}
		if len(idents) == 1 {
			if cc.List == nil {
				adData.Set(`adverts`, []struct{}{})
			} else {
				for _, rows := range cc.List {
					adData.Set(`adverts`, rows)
				}
			}
		} else {
			if cc.List == nil {
				adData.Set(`adverts`, map[string]struct{}{})
			} else {
				adData.Set(`adverts`, cc.List)
			}
		}
		if typ == `html` {
			cc.GenHTML()
		}
		data.SetData(adData)
	}
	return ctx.JSON(data)
}
