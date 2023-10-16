package customer

import (
	"github.com/webx-top/db"
)

func (f *Customer) IncrLoginFails() error {
	return f.OfficialCustomer.UpdateField(nil, `login_fails`, db.Raw(`login_fails+1`), `id`, f.Id)
}

func (f *Customer) ResetLoginFails() error {
	return f.OfficialCustomer.UpdateField(nil, `login_fails`, 0, `id`, f.Id)
}
