package resetpassword

import (
	"net/url"
	"sort"

	"github.com/coscms/webcore/library/common"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	"github.com/webx-top/echo"
)

type Validate func(c echo.Context, fieldName string, fieldValue string) error
type Send func(c echo.Context, m *modelCustomer.Customer, account string) error
type OnChange func(c echo.Context, m *modelCustomer.Customer) error
type GetAccount func(c echo.Context, m *modelCustomer.Customer) (string, error)

type RecvType struct {
	On             bool       // 开关
	Key            string     // Key
	Label          string     // 标签
	Placeholder    string     // 占位文本
	InputName      string     // 输入项名称
	Validate       Validate   `json:"-" xml:"-"` // 验证发送账号
	Send           Send       `json:"-" xml:"-"` // 发送
	OnChangeBefore OnChange   `json:"-" xml:"-"` // 更改密码前的处理
	OnChangeAfter  OnChange   `json:"-" xml:"-"` // 当密码成功更改后的处理
	GetAccount     GetAccount `json:"-" xml:"-"` // 收信账号
}

var recvTypes = map[string]*RecvType{
	`mobile`: &RecvType{
		On:            true,
		Key:           `mobile`,
		Label:         `手机`,
		Placeholder:   `手机号码`,
		InputName:     `手机号码`,
		Validate:      mobileValidate,
		Send:          mobileSend,
		GetAccount:    mobileGetAccount,
		OnChangeAfter: mobileOnChangeAfter,
	},
	`email`: &RecvType{
		On:            true,
		Key:           `email`,
		Label:         `E-mail`,
		Placeholder:   `E-mail地址`,
		InputName:     `E-mail地址`,
		Validate:      emailValidate,
		Send:          emailSend,
		GetAccount:    emailGetAccount,
		OnChangeAfter: emailOnChangeAfter,
	},
}

func Get(name string) *RecvType {
	if r, y := recvTypes[name]; y {
		return r
	}
	return nil
}

func Register(recvType *RecvType) {
	recvTypes[recvType.Key] = recvType
}

func Close(name string) {
	if v, y := recvTypes[name]; y {
		v.On = false
	}
}

func Open(name string) {
	if v, y := recvTypes[name]; y {
		v.On = true
	}
}

func Unregister(name string) {
	if _, y := recvTypes[name]; y {
		delete(recvTypes, name)
	}
}

func List() []*RecvType {
	var names []string
	for name, recvType := range recvTypes {
		if !recvType.On {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	rows := make([]*RecvType, len(names))
	for index, name := range names {
		rows[index] = recvTypes[name]
	}
	return rows
}

func GenResetPasswordURL(username string, typ string, account string) string {
	info := `name=` + url.QueryEscape(username) + `&type=` + url.QueryEscape(typ) + `&account=` + url.QueryEscape(account)
	encrypted := common.Crypto().Encode(info)
	resetPasswordURL := xMW.URLFor(`/forgot`) + `?vcode={code}&token=` + url.QueryEscape(encrypted)
	return resetPasswordURL
}
