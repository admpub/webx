package multidivicesignin

import (
	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/webx/application/library/xrole"
)

type MultideviceSignin struct {
	On         bool   `json:"on" xml:"on"`                 // 是否启用
	MaxDevices uint   `json:"maxDevices" xml:"maxDevices"` // 最大设备数量
	Unique     string `json:"unique" xml:"unique"`         // 需要限制唯一性的类型
}

func (a *MultideviceSignin) Combine(source interface{}) interface{} {
	src := source.(*MultideviceSignin)
	if src.On && !a.On {
		a.On = src.On
	}
	if src.MaxDevices > a.MaxDevices {
		a.MaxDevices = src.MaxDevices
	}
	if src.Unique != a.Unique {
		a.Unique = src.Unique
	}
	return a
}

const (
	UniqueDeviceID = `deviceID`
	UniquePlatform = `platform`
	UniqueScense   = `scense`
	BehaviorName   = `multi-device-signin`
)

func init() {
	xrole.Behaviors.Register(BehaviorName, `多设备登录`,
		perm.BehaviorOptFormHelpBlock(`配置多设备登录。on - 表示是否(true/false)启用多设备登录; maxDevices - 指定可以同时登录的最大设备数量; unique - 指定需要限制唯一性的类型(支持的值有:deviceID/platform/scense)`),
		perm.BehaviorOptValue(&MultideviceSignin{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &MultideviceSignin{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
}
