package customer

import (
	ngingdbschema "github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/db/lib/factory"
)

func init() {
	// - nging_file
	ngingdbschema.DBI.On(factory.EventCreated, func(m factory.Model, _ ...string) (err error) {
		fm := m.(*ngingdbschema.NgingFile)
		if fm.OwnerType == `customer` && fm.OwnerId > 0 {
			custM := customer.NewCustomer(fm.Context())
			err = custM.IncrFileSizeAndNum(fm.OwnerId, fm.Size, 1)
		}
		return
	}, `nging_file`)

	ngingdbschema.DBI.On(factory.EventDeleting, func(m factory.Model, _ ...string) (err error) {
		fm := m.(*ngingdbschema.NgingFile)
		if fm.OwnerType == `customer` && fm.OwnerId > 0 {
			custM := customer.NewCustomer(fm.Context())
			err = custM.DecrFileSizeAndNum(fm.OwnerId, fm.Size, 1)
		}
		return
	}, `nging_file`)
}
