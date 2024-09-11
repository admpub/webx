package customer

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/admpub/webx/application/dbschema"
)

func init() {
	// - official_customer
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCustomer)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Avatar
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_customer`, `avatar`).ListenDefault()
}
