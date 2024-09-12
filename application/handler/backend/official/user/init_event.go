package user

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	dbschemaNG "github.com/coscms/webcore/dbschema"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func init() {
	echo.OnCallback(`webx.customer.file.deleted`, func(v echo.Event) error {
		data := v.Context.Get(`data`).(*dbschemaNG.NgingFile)
		ownerID := v.Context.Uint64(`ownerID`)
		customerM := modelCustomer.NewCustomer(data.Context())
		err := customerM.Get(nil, db.Cond{`id`: ownerID})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return nil
			}
			return err
		}
		err = customerM.SafeDecrFileSizeAndNum(ownerID, data.Size, 1)
		if err != nil {
			_, _, err = customerM.RecountFile(ownerID)
		}
		return err
	})
}
