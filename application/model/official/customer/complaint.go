package customer

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
)

func NewComplaint(ctx echo.Context) *Complaint {
	m := &Complaint{
		OfficialCommonComplaint: dbschema.NewOfficialCommonComplaint(ctx),
	}
	return m
}

type Complaint struct {
	*dbschema.OfficialCommonComplaint
}

func (f *Complaint) check() error {
	if len(f.Content) == 0 && len(f.Type) == 0 {
		return f.Context().E(`请选择投诉类型`)
	}
	if f.CustomerId < 1 {
		return common.ErrUserNotLoggedIn
	}
	if f.TargetId < 1 && len(f.TargetIdent) == 0 {
		return f.Context().E(`投诉对象id无效`)
	}
	if len(f.TargetType) == 0 {
		return f.Context().E(`投诉对象类型无效`)
	}
	return nil
}

func (f *Complaint) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return
	}
	pk, err = f.OfficialCommonComplaint.Insert()
	if err != nil {
		return
	}
	err = ExecComplaintTargetFunc(f.TargetType, f.OfficialCommonComplaint)
	if err != nil {
		return
	}
	err = ExecComplaintTypeFunc(f.Type, f.OfficialCommonComplaint)
	return
}

func (f *Complaint) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonComplaint.Update(mw, args...)
}

func (f *Complaint) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCommonComplaint.Delete(mw, args...)
	return err
}

func (f *Complaint) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*ComplaintExt, error) {
	list := []*ComplaintExt{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...).Relation(`Customer`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			return sel.Columns(CusomterSafeFields...)
		})
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return list, err
	}
	for k, v := range list {
		v.TypeName = ComplaintTypes.Get(v.Type)
		item := ComplaintTargets.GetItem(v.TargetType)
		if item != nil {
			v.TargetTypeName = item.V
			if extraX, ok := item.X.(echo.H); ok {
				if urlFormat, ok := extraX.Get(`urlFormat`).(func(*dbschema.OfficialCommonComplaint) string); ok {
					v.SetURLFormat(urlFormat)
				}
			}
		}
		list[k] = v
	}
	return list, err
}
