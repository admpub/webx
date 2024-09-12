package invitation

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func CustomerList(ctx echo.Context) error {
	m := modelCustomer.NewInvitationCustomer(ctx)
	cond := db.NewCompounds()
	invitationID := ctx.Formx(`invitationId`).Uint()
	if invitationID > 0 {
		cond.AddKV(`invitation_id`, invitationID)
	}
	customerID := ctx.Formx(`customerId`).Uint()
	if customerID > 0 {
		cond.AddKV(`customer_id`, customerID)
	}
	list, err := m.ListCustomerWithCode(cond, `-id`)
	ctx.Set(`listData`, list)
	return ctx.Render(`official/customer/invitation/customer_list`, common.Err(ctx, err))
}
