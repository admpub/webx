package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/admpub/webx/application/dbschema"
)

func NewInvitationCustomer(ctx echo.Context) *InvitationCustomer {
	m := &InvitationCustomer{
		OfficialCustomerInvitationUsed: dbschema.NewOfficialCustomerInvitationUsed(ctx),
	}
	return m
}

type InvitationCustomer struct {
	*dbschema.OfficialCustomerInvitationUsed
}

func (f *InvitationCustomer) ListCustomer(cond *db.Compounds, orderby ...interface{}) ([]*InvitationCustomerExt, error) {
	list := []*InvitationCustomerExt{}
	err := f.listPage(&list, cond, orderby...)
	return list, err
}

func (f *InvitationCustomer) ListCustomerWithCode(cond *db.Compounds, orderby ...interface{}) ([]*InvitationCustomerWithCode, error) {
	list := []*InvitationCustomerWithCode{}
	err := f.listPage(&list, cond, orderby...)
	return list, err
}

func (f *InvitationCustomer) listPage(recv interface{}, cond *db.Compounds, orderby ...interface{}) error {
	_, err := common.NewLister(f, recv, func(r db.Result) db.Result {
		return r.OrderBy(orderby...).Relation(`Customer`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			return sel.Columns(CusomterSafeFields...)
		})
	}, cond.And()).Paging(f.Context())
	return err
}
