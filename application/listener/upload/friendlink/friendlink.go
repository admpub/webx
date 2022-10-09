package friendlink

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/v4/application/library/fileupdater/listener"
	"github.com/admpub/nging/v4/application/registry/upload/thumb"
	"github.com/admpub/webx/application/dbschema"
)

var FriendlinkLogoThumbnail = thumb.Size{
	Width:  88 * 2,
	Height: 31 * 2,
}

func init() {
	thumb.Registry.Set(`friendlink`, FriendlinkLogoThumbnail)
	//友情链接
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonFriendlink)
		tableID = fmt.Sprint(fm.Id)
		content = fm.LogoOriginal
		property = listener.NewPropertyWith(
			fm,
			db.Cond{`id`: fm.Id},
			listener.FieldValueWith(`image`, thumb.DefaultSize.ThumbValue()),
		)
		return
	}, false).SetTable(`official_common_friendlink`, `logo_original`, `logo`).ListenDefault()
}
