package customer

import (
	"github.com/admpub/log"
	"github.com/webx-top/db"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/sessionguard"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/middleware/sessdata"
)

// SetSession 记录登录信息
func (f *Customer) SetSession(customers ...*dbschema.OfficialCustomer) {
	customerCopy := f.ClearPasswordData(customers...)
	f.Context().Session().Set(`customer`, &customerCopy)
}

// UnsetSession 退出登录
func (f *Customer) UnsetSession() error {
	err := FireSignOut(f.OfficialCustomer)
	f.Context().Session().Delete(`customer`)
	sessdata.ClearPermissionCache(f.OfficialCustomer.Id)
	return err
}

func (f *Customer) VerifySession(customers ...*dbschema.OfficialCustomer) error {
	var customer *dbschema.OfficialCustomer
	if len(customers) > 0 {
		customer = customers[0]
	} else {
		customer, _ = f.Context().Session().Get(`customer`).(*dbschema.OfficialCustomer)
	}
	if customer == nil {
		return common.ErrUserNotLoggedIn
	}
	detail, err := f.GetDetail(db.Cond{`id`: customer.Id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		f.UnsetSession()
		return common.ErrUserNotFound
	}
	if detail.OfficialCustomer.SessionId != f.Context().Session().ID() {
		f.UnsetSession()
		return common.ErrUserNotLoggedIn
	}
	if !sessionguard.Validate(f.Context(), ``, `customer`, detail.Id) {
		log.Warn(f.Context().T(`客户“%s”的会话环境发生改变，需要重新登录`, detail.Name))
		f.UnsetSession()
		return common.ErrUserNotLoggedIn
	}
	if detail.OfficialCustomer.Updated != customer.Updated {
		f.SetSession(detail.OfficialCustomer)
		f.Context().Internal().Set(`customer`, detail.OfficialCustomer)
	}
	if detail.OfficialCustomer.RoleIds != customer.RoleIds {
		sessdata.ClearPermissionCache(detail.OfficialCustomer.Id)
	}

	safeCustomer := f.ClearPasswordData(detail.OfficialCustomer)
	detail.OfficialCustomer = &safeCustomer
	f.Context().Internal().Set(`customerDetail`, detail)
	return nil
}
