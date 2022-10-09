package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func (f *Customer) IncrLoginFails() error {
	return f.OfficialCustomer.UpdateField(nil, `login_fails`, db.Raw(`login_fails+1`), `id`, f.Id)
}

func (f *Customer) ResetLoginFails() error {
	return f.OfficialCustomer.UpdateField(nil, `login_fails`, 0, `id`, f.Id)
}

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
