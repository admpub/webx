package binding

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webfront/dbschema"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

// BindOAuthAccount 第三方平台信息
type BindOAuthAccount struct {
	IconClass   string //图标class属性值
	IconImage   string //图片图标
	WrapClass   string
	Provider    string                            //平台标识
	Title       string                            //平台账号名称
	Description string                            //绑定说明
	Binded      bool                              //当前账号是否已经绑定
	OAuthUsers  []*dbschema.OfficialCustomerOauth //已绑定用户列表
}

// AccountBind 账号绑定操作信息
type AccountBind struct {
	Type       string                                            //类型标识(比如，email)
	Name       string                                            //类型名(中文。比如，邮箱)
	ObjectName string                                            //发送物件名称(比如，向邮箱发邮件，向手机发短信。这里的"邮件","短信"即是发送物件)
	Verifier   func(echo.Context, *modelCustomer.Customer) error //验证收到的验证码逻辑
	Sender     func(echo.Context, *modelCustomer.Customer) error //发信逻辑
}

var (
	oAuthProviders = []*BindOAuthAccount{}
	accountBinders = map[string]*AccountBind{}
	NoopBinder     = func(_ echo.Context, _ *modelCustomer.Customer) error {
		return nil
	}
)

func Register(a *AccountBind) {
	if a.Sender == nil {
		a.Sender = NoopBinder
	}
	if a.Verifier == nil {
		a.Verifier = NoopBinder
	}
	accountBinders[a.Type] = a
}

func Get(typ string) *AccountBind {
	a, _ := accountBinders[typ]
	return a
}

func All() map[string]*AccountBind {
	return accountBinders
}
