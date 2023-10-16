package customer

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/webx/application/dbschema"
	multidivicesignin "github.com/admpub/webx/application/library/multidevicesignin"
)

// NewDevice 客户登录设备
func NewDevice(ctx echo.Context) *Device {
	m := &Device{
		OfficialCustomerDevice: dbschema.NewOfficialCustomerDevice(ctx),
	}
	return m
}

type Device struct {
	*dbschema.OfficialCustomerDevice
}

func (f *Device) Exists(customerID uint64, scense, platform, deviceNo string) (bool, error) {
	err := f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`scense`: scense},
		db.Cond{`platform`: platform},
		db.Cond{`device_no`: deviceNo},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (f *Device) ExistsOther(customerID uint64, scense, platform, deviceNo string, excludeIDs uint64) (bool, error) {
	err := f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`scense`: scense},
		db.Cond{`platform`: platform},
		db.Cond{`device_no`: deviceNo},
		db.Cond{`id`: db.NotEq(excludeIDs)},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (f *Device) SetDefaults() {
	if len(f.SessionId) == 0 {
		f.SessionId = f.Context().Session().MustID()
	}
	if len(f.Platform) == 0 {
		f.Platform = DefaultDevicePlatform // pc/android/ios/micro-program
	}
	if len(f.Scense) == 0 {
		f.Scense = DefaultDeviceScense // app/web
	}
	if len(f.DeviceNo) == 0 {
		f.DeviceNo = `none`
	}
}

func (f *Device) check() error {
	f.SetDefaults()
	var (
		exists bool
		err    error
	)
	if f.Id <= 0 {
		exists, err = f.Exists(f.CustomerId, f.Scense, f.Platform, f.DeviceNo)
	} else {
		exists, err = f.ExistsOther(f.CustomerId, f.Scense, f.Platform, f.DeviceNo, f.Id)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().NewError(code.DataAlreadyExists, `数据已经存在: [ customer_id: %d ; scense: %s ; platform: %s ; device_no: %s ]`, f.CustomerId, f.Scense, f.Platform, f.DeviceNo)
	}
	return err
}

func (f *Device) Add() (pk interface{}, err error) {
	err = f.check()
	if err != nil {
		return
	}
	f.Updated = uint(time.Now().Unix())
	pk, err = f.OfficialCustomerDevice.Insert()
	return
}

func (f *Device) Upsert() (pk interface{}, err error) {
	f.SetDefaults()
	row := dbschema.NewOfficialCustomerDevice(f.Context())
	cond := db.And(
		db.Cond{`customer_id`: f.CustomerId},
		db.Cond{`scense`: f.Scense},
		db.Cond{`platform`: f.Platform},
		db.Cond{`device_no`: f.DeviceNo},
	)
	err = row.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `created`)
	}, cond)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		f.Updated = uint(time.Now().Unix())
		return f.OfficialCustomerDevice.Insert()
	}
	f.OfficialCustomerDevice.Id = row.Id
	f.OfficialCustomerDevice.Created = row.Created
	err = f.OfficialCustomerDevice.Update(nil, `id`, row.Id)
	return
}

func (f *Device) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCustomerDevice.Update(mw, args...)
}

func (f *Device) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCustomerDevice.Delete(mw, args...)
	return err
}

func (f *Device) SignOut(mw func(db.Result) db.Result, args ...interface{}) error {
	cnt, err := f.OfficialCustomerDevice.ListByOffset(nil, mw, 0, 500, args...)
	if err != nil {
		return err
	}
	for _, row := range f.Objects() {
		if len(row.SessionId) > 0 {
			if err = f.Context().Session().RemoveID(row.SessionId); err != nil {
				return err
			}
		}
		if err = f.OfficialCustomerDevice.Delete(nil, `id`, row.Id); err != nil {
			return err
		}
	}
	if cnt() > 0 {
		return f.SignOut(mw, args...)
	}
	return err
}

func (f *Device) CleanCustomer(customer *dbschema.OfficialCustomer, options ...CustomerOption) (err error) {
	multideviceSignin, ok := CustomerRolePermissionForBehavior(f.Context(), multidivicesignin.BehaviorName, customer).(*multidivicesignin.MultideviceSignin)
	cond := db.NewCompounds()
	cond.Add(db.Cond{`customer_id`: customer.Id})
	co := NewCustomerOptions(customer)
	if !ok {
		goto END
	}
	for _, option := range options {
		option(co)
	}
	f.SetOptions(co)
	f.SetDefaults()
	switch multideviceSignin.Unique {
	case multidivicesignin.UniquePlatform:
		cond.Add(db.Cond{`scense`: f.Scense})
		cond.Add(db.Cond{`platform`: f.Platform})
	case multidivicesignin.UniqueScense:
		cond.Add(db.Cond{`scense`: f.Scense})
	case multidivicesignin.UniqueDeviceID:
		fallthrough
	default:
		cond.Add(db.Cond{`scense`: f.Scense})
		cond.Add(db.Cond{`platform`: f.Platform})
		cond.Add(db.Cond{`device_no`: f.DeviceNo})
	}

END:
	if f.Id > 0 {
		cond.Add(db.Cond{`id`: db.NotEq(f.Id)})
	}
	err = f.SignOut(nil, cond.And())
	if err == nil {
		err = f.CleanExpired()
	}
	if err == nil {
		err = f.CleanExceedLimit(customer.Id, multideviceSignin)
	}
	return
}

func (f *Device) CleanExceedLimit(customerID uint64, multideviceSignin *multidivicesignin.MultideviceSignin) error {
	if multideviceSignin == nil || multideviceSignin.MaxDevices <= 0 {
		return nil
	}
	oldRows := []*dbschema.OfficialCustomerDevice{}
	oldCond := db.NewCompounds()
	oldCond.Add(db.Cond{`customer_id`: customerID})
	_, err := f.OfficialCustomerDevice.ListByOffset(&oldRows, func(r db.Result) db.Result {
		return r.Select(`id`, `updated`).OrderBy(`-updated`, `-id`)
	}, 0, int(multideviceSignin.MaxDevices), oldCond.And())
	if err != nil {
		return err
	}
	if len(oldRows) == int(multideviceSignin.MaxDevices) {
		excludeIDs := make([]uint64, len(oldRows))
		for i, r := range oldRows {
			excludeIDs[i] = r.Id
		}
		if f.Id > 0 && !com.InUint64Slice(f.Id, excludeIDs) {
			excludeIDs = append(excludeIDs, f.Id)
		}
		lastRow := oldRows[len(oldRows)-1]
		if lastRow.Updated > 0 {
			oldCond.Add(db.Cond{`updated`: db.Lte(lastRow.Updated)})
		}
		oldCond.Add(db.Cond{`id`: db.NotIn(excludeIDs)})
		err = f.SignOut(nil, oldCond.And())
	}
	return err
}

func (f *Device) CleanExpired() (err error) {
	return f.SignOut(nil, db.And(
		db.Cond{`session_id`: db.NotEq(``)},
		db.Cond{`expired`: db.NotEq(0)},
		db.Cond{`expired`: db.Lte(time.Now().Unix())},
	))
}

func (f *Device) Kick(customerID uint64) (err error) {
	return f.SignOut(nil, db.Cond{`customer_id`: customerID})
}

func (f *Device) SetOptions(options *CustomerOptions) *Device {
	f.DeviceNo = options.DeviceNo
	f.Platform = options.Platform
	f.Scense = options.Scense
	if options.MaxAge > 0 {
		f.Expired = uint(time.Now().Add(options.MaxAge).Unix())
	}
	return f
}
