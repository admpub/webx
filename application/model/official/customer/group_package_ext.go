package customer

import "github.com/webx-top/echo"

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
