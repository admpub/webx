package customer

import (
	"errors"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/webx/application/dbschema"
)

func NewInvitation(ctx echo.Context) *Invitation {
	m := &Invitation{
		OfficialCustomerInvitation: dbschema.NewOfficialCustomerInvitation(ctx),
	}
	return m
}

type Invitation struct {
	*dbschema.OfficialCustomerInvitation
}

func (f *Invitation) FindCode(invitationCode string) error {
	err := f.Get(nil, `code`, invitationCode)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = f.Context().E(`邀请码不存在`)
		}
		return err
	}
	if f.Disabled == `Y` {
		return f.Context().E(`邀请码无效`)
	}
	if f.AllowNum <= f.UsedNum {
		if f.AllowNum == 1 {
			err = f.Context().E(`邀请码“%s”已经使用过了`, invitationCode)
		} else {
			err = f.Context().E(`邀请码“%s”已经达到最大使用次数`, invitationCode)
		}
	}
	now := uint(time.Now().Unix())
	if f.Start > now {
		if f.End > 0 {
			err = errors.New(f.Context().T(`该邀请码只能在“%s - %s”这段时间内使用`,
				tplfunc.TsToDate(`2006/01/02 15:04:05`, f.Start),
				tplfunc.TsToDate(`2006/01/02 15:04:05`, f.End),
			))
		} else {
			err = errors.New(f.Context().T(`该邀请码只能在“%s”之后使用`,
				tplfunc.TsToDate(`2006/01/02 15:04:05`, f.Start),
			))
		}
		return err
	}
	if f.End > 0 && f.End < now {
		err = f.Context().E(`该邀请码已过期`)
		return err
	}
	return err
}

func (f *Invitation) check() error {
	if len(f.Code) == 0 {
		return f.Context().E(`邀请码不能为空`)
	}
	exists, err := f.Exists(f.Code, f.Id)
	if err != nil {
		return err
	}
	if exists {
		return f.Context().E(`邀请码已经存在`)
	}
	return nil
}

func (f *Invitation) Exists(code string, excludeIDs ...uint) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV(`code`, code)
	if len(excludeIDs) > 0 && excludeIDs[0] > 0 {
		cond.Add(db.Cond{`id`: db.NotEq(excludeIDs[0])})
	}
	return f.OfficialCustomerInvitation.Exists(nil, cond.And())
}

func (f *Invitation) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCustomerInvitation.Insert()
}

func (f *Invitation) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCustomerInvitation.Update(mw, args...)
}

func (f *Invitation) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCustomerInvitation.Delete(mw, args...)
	return err
}

func (f *Invitation) UseCode(invitationID uint, customer *dbschema.OfficialCustomer) error {
	err := f.UpdateField(nil, `used_num`, db.Raw(`used_num+1`), `id`, invitationID)
	if err != nil {
		return err
	}
	m := dbschema.NewOfficialCustomerInvitationUsed(f.Context())
	m.CustomerId = customer.Id
	m.InvitationId = invitationID
	m.LevelId = customer.LevelId
	m.AgentLevelId = customer.AgentLevel
	m.RoleIds = customer.RoleIds
	_, err = m.Insert()
	return err
}
