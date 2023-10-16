package user

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	dbschemaNG "github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/webx/application/dbschema"
)

func init() {
	echo.OnCallback(`webx.customer.file.deleted`, func(v echo.Event) error {
		data := v.Context.Get(`data`).(*dbschemaNG.NgingFile)
		ownerID := v.Context.Uint64(`ownerID`)
		customerM := dbschema.NewOfficialCustomer(data.Context())
		err := customerM.Get(nil, db.Cond{`id`: ownerID})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return nil
			}
			return err
		}
		err = customerM.UpdateFields(nil, map[string]interface{}{
			`file_size`: db.Raw(`file_size-` + fmt.Sprintf(`%d`, data.Size)),
			`file_num`:  db.Raw(`file_num-1`),
		}, db.And(
			db.Cond{`id`: ownerID},
			db.Cond{`file_size`: db.Gte(data.Size)},
			db.Cond{`file_num`: db.Gt(0)},
		))
		if err != nil {
			fileM := dbschemaNG.NewNgingFile(data.Context())
			recv := echo.H{}
			err = fileM.NewParam().SetMW(func(r db.Result) db.Result {
				return r.Select(db.Raw(`SUM(size) AS c`), db.Raw(`COUNT(1) AS n`))
			}).SetRecv(&recv).SetArgs(db.And(
				db.Cond{`owner_type`: `user`},
				db.Cond{`owner_id`: ownerID},
			)).One()
			if err != nil {
				return err
			}
			totalNum := recv.Uint64(`n`)
			totalSize := recv.Uint64(`c`)
			err = customerM.UpdateFields(nil, map[string]interface{}{
				`file_size`: totalSize,
				`file_num`:  totalNum,
			}, db.Cond{`id`: ownerID})
		}
		return err
	})
}
