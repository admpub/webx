package xcommon

import (
	"errors"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/webx/application/library/xrole"
	"github.com/webx-top/echo"
)

type CustomerAddedCounter interface {
	CustomerTodayCount(uid interface{}) (int64, error)
	CustomerPendingCount(uid interface{}) (int64, error)
	CustomerPendingTodayCount(uid interface{}) (int64, error)
}

type ConfigCustomerAdd struct {
	MaxPerDay        int64 `json:"maxPerDay"`
	MaxPending       int64 `json:"maxPending"`
	MaxPendingPerDay int64 `json:"maxPendingPerDay"`
}

func (c *ConfigCustomerAdd) Combine(source interface{}) interface{} {
	src := source.(*ConfigCustomerAdd)
	if src.MaxPerDay > c.MaxPerDay {
		c.MaxPerDay = src.MaxPerDay
	}
	if src.MaxPending > c.MaxPending {
		c.MaxPending = src.MaxPending
	}
	if src.MaxPendingPerDay > c.MaxPendingPerDay {
		c.MaxPendingPerDay = src.MaxPendingPerDay
	}
	return c
}

var (
	ErrCustomerAddClosed           = errors.New(`已关闭`)
	ErrCustomerAddMaxPerDay        = errors.New(`已达到今日最大数量`)
	ErrCustomerAddMaxPending       = errors.New(`待审核数量已达上限`)
	ErrCustomerAddMaxPendingPerDay = errors.New(`待审核数量已达今日上限`)
)

func CheckRoleCustomerAdd(ctx echo.Context, permission *xrole.RolePermission, behaviorName string, customerID interface{}, counter CustomerAddedCounter) error {
	if permission == nil {
		return CheckGlobalCustomerAdd(ctx, customerID, behaviorName, counter)
	}
	bev, ok := permission.Get(ctx, xrole.CustomerRolePermissionTypeBehavior).(perm.BehaviorPerms)
	if !ok {
		return CheckGlobalCustomerAdd(ctx, customerID, behaviorName, counter)
	}
	roleCfg, _ := bev.Get(behaviorName).Value.(*ConfigCustomerAdd)
	if roleCfg == nil {
		return CheckGlobalCustomerAdd(ctx, customerID, behaviorName, counter)
	}
	if roleCfg.MaxPerDay <= 0 {
		return ErrCustomerAddClosed
	}
	todayCount, err := counter.CustomerTodayCount(customerID)
	if err != nil {
		return err
	}
	if todayCount >= roleCfg.MaxPerDay {
		return ErrCustomerAddMaxPerDay
	}
	if roleCfg.MaxPending > 0 {
		pendingCount, err := counter.CustomerPendingCount(customerID)
		if err != nil {
			return err
		}
		if pendingCount >= roleCfg.MaxPending {
			return ErrCustomerAddMaxPending
		}
	}
	if roleCfg.MaxPendingPerDay > 0 {
		pendingPerDayCount, err := counter.CustomerPendingTodayCount(customerID)
		if err != nil {
			return err
		}
		if pendingPerDayCount >= roleCfg.MaxPendingPerDay {
			return ErrCustomerAddMaxPendingPerDay
		}
	}
	return err
}

func CheckGlobalCustomerAdd(ctx echo.Context, customerID interface{}, typ string, counter CustomerAddedCounter) error {
	cmtFrequency := common.Setting(`frequency`, typ)
	var maxPerDay int64
	if len(cmtFrequency.String(`maxPerDay`)) > 0 {
		maxPerDay = cmtFrequency.Int64(`maxPerDay`)
	} else {
		maxPerDay = 100
	}
	if maxPerDay <= 0 {
		return ErrCustomerAddClosed
	}
	todayCount, err := counter.CustomerTodayCount(customerID)
	if err != nil {
		return err
	}
	if todayCount >= maxPerDay {
		return ErrCustomerAddMaxPerDay
	}
	var maxPending int64
	if len(cmtFrequency.String(`maxPending`)) > 0 {
		maxPending = cmtFrequency.Int64(`maxPending`)
	}
	if maxPending > 0 {
		pendingCount, err := counter.CustomerPendingCount(customerID)
		if err != nil {
			return err
		}
		if pendingCount >= maxPending {
			return ErrCustomerAddMaxPending
		}
	}
	var maxPendingPerDay int64
	if len(cmtFrequency.String(`maxPendingPerDay`)) > 0 {
		maxPendingPerDay = cmtFrequency.Int64(`maxPendingPerDay`)
	}
	if maxPendingPerDay > 0 {
		pendingPerDayCount, err := counter.CustomerPendingCount(customerID)
		if err != nil {
			return err
		}
		if pendingPerDayCount >= maxPendingPerDay {
			return ErrCustomerAddMaxPendingPerDay
		}
	}
	return err
}
