package customer

import "github.com/admpub/webx/application/dbschema"

type WalletExt struct {
	*dbschema.OfficialCustomerWallet
	Customer      *dbschema.OfficialCustomer `db:"-,relation=id:customer_id"`
	AssetTypeName string                     `db:"-"`
}

type WalletFlowExt struct {
	*dbschema.OfficialCustomerWalletFlow
	Customer      *dbschema.OfficialCustomer `db:"-,relation=id:customer_id"`
	SrcCustomer   *dbschema.OfficialCustomer `db:"-,relation=id:source_customer|gtZero"`
	AssetTypeName string                     `db:"-"`
}
