package customer

import (
	"github.com/webx-top/echo"
)

const (
	GroupPackageTimeDay     = `day`
	GroupPackageTimeWeek    = `week`
	GroupPackageTimeMonth   = `month`
	GroupPackageTimeYear    = `year`
	GroupPackageTimeForever = `forever`
)

var GroupPackageTimeUnits = echo.NewKVData().
	Add(GroupPackageTimeDay, `天`).
	Add(GroupPackageTimeWeek, `周`).
	Add(GroupPackageTimeMonth, `月`).
	Add(GroupPackageTimeYear, `年`).
	Add(GroupPackageTimeForever, `永久`)

func GroupPackageTimeUnitSuffix(c echo.Context, n uint, unit string) string {
	if unit == GroupPackageTimeMonth {
		if n > 1 {
			return c.T(`/ %d个月`, n)
		}
		return c.T(`/ ` + GroupPackageTimeUnits.Get(unit))
	}
	if unit == GroupPackageTimeForever {
		return c.T(GroupPackageTimeUnits.Get(unit))
	}
	if n > 1 {
		return c.T(`/ %d`+GroupPackageTimeUnits.Get(unit), n)
	}
	return c.T(`/ ` + GroupPackageTimeUnits.Get(unit))
}
