package category

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/admpub/webx/application/dbschema"
)

func init() {

	// - official_common_category
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonCategory)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Cover
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_common_category`, `cover`).ListenDefault()
}
