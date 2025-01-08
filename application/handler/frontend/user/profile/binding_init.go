package profile

import (
	"github.com/admpub/webx/application/handler/frontend/user/binding"
	"github.com/webx-top/echo"
)

func init() {
	binding.Register(&binding.AccountBind{
		Type:       `email`,
		Name:       echo.T(`邮箱`),
		ObjectName: `邮件`,
		Verifier:   bindingEmailVerify,
		Sender:     bindingEmailSend,
	})
	binding.Register(&binding.AccountBind{
		Type:       `mobile`,
		Name:       echo.T(`手机`),
		ObjectName: `短信`,
		Verifier:   bindingMobileVerify,
		Sender:     bindingMobileSend,
	})
	binding.Register(&binding.AccountBind{
		Type:       `oauth`,
		Name:       echo.T(`第三方`),
		ObjectName: `账号`,
		Verifier:   bindingOAuthGet,
		Sender:     bindingOAuth,
	})
}
