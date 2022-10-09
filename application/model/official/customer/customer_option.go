package customer

import (
	"time"

	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

func NewCustomerOptions(customerM *dbschema.OfficialCustomer) *CustomerOptions {
	var customer *dbschema.OfficialCustomer
	if customerM == nil {
		customer = dbschema.NewOfficialCustomer(nil)
	} else {
		_customer := ClearPasswordData(customerM)
		customer = &_customer
	}
	return &CustomerOptions{OfficialCustomer: customer}
}

type CustomerOptions struct {
	*dbschema.OfficialCustomer
	MaxAge     time.Duration // 登录状态有效时长
	SignInType string        // 登录方式
	Scense     string        // 场景
	Platform   string        // 系统平台
	DeviceNo   string        // 设备编号
}

type CustomerOption func(*CustomerOptions)

func CustomerName(name string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Name = name
	}
}
func CustomerPassword(password string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Password = password
	}
}
func CustomerMobile(mobile string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Mobile = mobile
	}
}
func CustomerEmail(email string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Email = email
	}
}
func CustomerMaxAgeSeconds(maxAgeSeconds int) CustomerOption {
	return func(c *CustomerOptions) {
		c.MaxAge = time.Duration(maxAgeSeconds) * time.Second
	}
}
func CustomerMaxAge(maxAge time.Duration) CustomerOption {
	return func(c *CustomerOptions) {
		c.MaxAge = maxAge
	}
}
func CustomerSignInType(signInType string) CustomerOption {
	return func(c *CustomerOptions) {
		c.SignInType = signInType
	}
}
func CustomerScense(scense string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Scense = scense
	}
}
func CustomerPlatform(platform string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Platform = platform
	}
}
func CustomerDeviceNo(deviceNo string) CustomerOption {
	return func(c *CustomerOptions) {
		c.DeviceNo = deviceNo
	}
}
func GenerateOptionsFromHeader(c echo.Context, maxAge ...int) []CustomerOption {
	co := []CustomerOption{
		CustomerPlatform(c.Header(`X-Platform`)),
		CustomerScense(c.Header(`X-Scense`)),
		CustomerDeviceNo(c.Header(`X-Device-Id`)),
	}
	if len(maxAge) > 0 {
		co = append(co, CustomerMaxAgeSeconds(maxAge[0]))
	}
	return co
}
