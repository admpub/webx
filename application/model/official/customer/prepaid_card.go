package customer

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/errors"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
)

var ErrCardNumberAlreadyExists = errors.New(`Card number already exists`)

func NewPrepaidCard(ctx echo.Context) *PrepaidCard {
	m := &PrepaidCard{
		OfficialCustomerPrepaidCard: dbschema.NewOfficialCustomerPrepaidCard(ctx),
	}
	return m
}

type PrepaidCard struct {
	*dbschema.OfficialCustomerPrepaidCard
}

func (f *PrepaidCard) check() error {
	if f.Uid < 1 {
		return common.ErrUserNotLoggedIn
	}
	if f.Amount < 1 {
		return f.Context().E(`面值无效，必须大于0`)
	}
	if f.SalePrice <= 0 {
		return f.Context().E(`售价无效，必须大于0`)
	}
	exists, err := f.Exists(f.Number, f.Id)
	if err != nil {
		return err
	}
	if exists {
		return ErrCardNumberAlreadyExists
	}
	return nil
}

func (f *PrepaidCard) UseCard(customerID uint64, number string, password string) error {
	err := f.Get(nil, db.Cond{`number`: number})
	if err != nil {
		if err == db.ErrNoMoreRows {
			return f.Context().E(`充值卡不存在`)
		}
		return err
	}
	if f.Disabled == `Y` {
		return f.Context().E(`充值卡无效`)
	}
	if f.Password != password {
		return f.Context().E(`充值卡密码错误`)
	}
	if f.Used > 0 {
		return f.Context().E(`充值卡“%v”已经使用过了`, f.Number)
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
	kvset := echo.H{
		`used`:        uint(time.Now().Unix()),
		`customer_id`: customerID,
	}
	return f.UpdateFields(nil, kvset, `id`, f.Id)
}

// BatchGenerate 批量生成
func (f *PrepaidCard) BatchGenerate(uid uint, count int, amount uint, salePrice float64, start uint, end uint, bgImage string) error {
	for i := 0; i < count; i++ {
		f.Reset()
		f.Uid = uid
		f.Amount = amount
		f.SalePrice = salePrice
		f.Start = start
		f.End = end
		f.BgImage = bgImage
		_, err := f.Add()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *PrepaidCard) Add() (pk interface{}, err error) {
	f.Number, err = common.UniqueID()
	if err != nil {
		return nil, err
	}
	cardNumber := f.Number
	var loopTimes int
	loopMaxTimes := 10

LOOP:
	if err = f.check(); err != nil {
		if err == ErrCardNumberAlreadyExists {
			if loopTimes >= loopMaxTimes {
				return nil, errors.WithMessage(err, f.Number)
			}
			f.Number = cardNumber + com.RandomNumeric(2)
			loopTimes++
			goto LOOP
		}
		return nil, err
	}
	f.Password = com.RandomNumeric(8)
	return f.OfficialCustomerPrepaidCard.Insert()
}

func (f *PrepaidCard) Exists(number string, excludeIDs ...uint64) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV(`number`, number)
	if len(excludeIDs) > 0 && excludeIDs[0] > 0 {
		cond.Add(db.Cond{`id`: db.NotEq(excludeIDs[0])})
	}
	return f.OfficialCustomerPrepaidCard.Exists(nil, cond.And())
}

func (f *PrepaidCard) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCustomerPrepaidCard.Update(mw, args...)
}

func (f *PrepaidCard) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCustomerPrepaidCard.Delete(mw, args...)
	return err
}
