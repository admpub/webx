package customer

import "github.com/webx-top/echo"

var (
	DevicePlatforms       = echo.NewKVData()
	DeviceScenses         = echo.NewKVData()
	DefaultDevicePlatform = `pc`
	DefaultDeviceScense   = `web`
)

func init() {
	DevicePlatforms.Add(`ios`, `iOS手机`)
	DevicePlatforms.Add(`android`, `安卓手机`)
	DevicePlatforms.Add(`pc`, `个人电脑`)
	DevicePlatforms.Add(`micro-program`, `小程序`)

	DeviceScenses.Add(`app`, `App应用`)
	DeviceScenses.Add(`web`, `网页`)
}
