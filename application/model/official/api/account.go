package api

import (
	"encoding/json"

	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewAccount(ctx echo.Context) *Account {
	m := &Account{
		OfficialCommonApiAccount: dbschema.NewOfficialCommonApiAccount(ctx),
	}
	return m
}

type Account struct {
	*dbschema.OfficialCommonApiAccount
}

func (f *Account) check() error {
	if len(f.AppId) == 0 {
		return f.Context().NewError(code.InvalidParameter, `appId is required`).SetZone(`appId`)
	}
	if len(f.OwnerType) == 0 {
		return f.Context().NewError(code.InvalidParameter, `ownerType is required`).SetZone(`ownerType`)
	}
	if f.OwnerId < 1 {
		return f.Context().NewError(code.InvalidParameter, `ownerId is required`).SetZone(`ownerId`)
	}
	return nil
}

func (f *Account) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	if err := f.Exists(f.AppId); err != nil {
		return nil, err
	}
	return f.OfficialCommonApiAccount.Insert()
}

func (f *Account) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	if err := f.ExistsOther(f.AppId, f.Id); err != nil {
		return err
	}
	return f.OfficialCommonApiAccount.Update(mw, args...)
}

func (f *Account) Exists(appId string) error {
	exists, err := f.OfficialCommonApiAccount.Exists(nil, db.Cond{`app_id`: appId})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`AppId“%s”已经存在`, appId)
	}
	return err
}

func (f *Account) ExistsOther(appId string, id uint64) error {
	exists, err := f.OfficialCommonApiAccount.Exists(nil, db.And(
		db.Cond{`app_id`: appId},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`AppId“%s”已经存在`, appId)
	}
	return err
}

func (f *Account) ParseExtra() echo.H {
	r := echo.H{}
	if len(f.Extra) > 0 {
		json.Unmarshal([]byte(f.Extra), &r)
	}
	return r
}

func (f *Account) GetByAppID(appID string) (*dbschema.OfficialCommonApiAccount, error) {
	err := f.Get(nil, `app_id`, appID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = f.Context().NewError(code.DataNotFound, `API账号没找到`)
		}
	}
	return f.OfficialCommonApiAccount, err
}

func (f *Account) OauthLogin() echo.H {
	r := echo.H{}
	// TODO:
	return r
}
