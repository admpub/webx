package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	dbschemaNG "github.com/coscms/webcore/dbschema"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/xrole/xroleupload"
)

func (f *Customer) IncrFileSizeAndNum(customerID uint64, fileSize uint64, fileNum uint64) error {
	return f.OfficialCustomer.UpdateFields(nil, echo.H{
		`file_size`: db.Raw(`file_size+` + param.AsString(fileSize)),
		`file_num`:  db.Raw(`file_num+` + param.AsString(fileNum)),
	}, `id`, customerID)
}

func (f *Customer) DecrFileSizeAndNum(customerID uint64, fileSize uint64, fileNum uint64) error {
	return f.OfficialCustomer.UpdateFields(nil, echo.H{
		`file_size`: db.Raw(`file_size-` + param.AsString(fileSize)),
		`file_num`:  db.Raw(`file_num-` + param.AsString(fileNum)),
	}, `id`, customerID)
}

func (f *Customer) SafeDecrFileSizeAndNum(customerID uint64, fileSize uint64, fileNum uint64) error {
	return f.UpdateFields(nil, map[string]interface{}{
		`file_size`: db.Raw(`file_size-` + param.AsString(fileSize)),
		`file_num`:  db.Raw(`file_num-` + param.AsString(fileNum)),
	}, db.And(
		db.Cond{`id`: customerID},
		db.Cond{`file_size`: db.Gte(fileSize)},
		db.Cond{`file_num`: db.Gte(fileNum)},
	))
}

// RecountFile 重新统计客户上传的文件数量和尺寸
func (f *Customer) RecountFile(customerId ...uint64) (totalNum uint64, totalSize uint64, err error) {
	ownerID := f.Id
	if len(customerId) > 0 && customerId[0] > 0 {
		ownerID = customerId[0]
	}
	fileM := dbschemaNG.NewNgingFile(f.Context())
	recv := echo.H{}
	err = fileM.NewParam().SetMW(func(r db.Result) db.Result {
		return r.Select(db.Raw(`SUM(size) AS c`), db.Raw(`COUNT(1) AS n`))
	}).SetRecv(&recv).SetArgs(db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: ownerID},
	)).One()
	if err != nil {
		return
	}
	totalNum = recv.Uint64(`n`)
	totalSize = recv.Uint64(`c`)
	err = f.UpdateFields(nil, map[string]interface{}{
		`file_size`: totalSize,
		`file_num`:  totalNum,
	}, db.Cond{`id`: ownerID})
	return
}

func (f *Customer) GetUploadConfig(customer *dbschema.OfficialCustomer) *xroleupload.CustomerUpload {
	cfg, _ := CustomerRolePermissionForBehavior(f.Context(), xroleupload.BehaviorName, customer).(*xroleupload.CustomerUpload)
	return cfg
}
