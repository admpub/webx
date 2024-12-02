package index

import (
	"github.com/admpub/nging/v5/application/handler/manager"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.Register(func(g echo.RouteRegister) {
		g.Route(`GET,HEAD`, uploadLibrary.UploadURLPath+`:subdir/*`, manager.File)
		g.Route(`GET,POST`, `/`, Index).SetName(`root`)
		g.Route(`GET,POST`, `/index`, Index).SetName(`index`)
		g.Route(`GET,POST`, `/search`, Search).SetName(`search`)
		g.Route(`GET,POST`, `/sign_up`, SignUp).SetName(`sign_up`)
		g.Route(`GET,POST`, `/sign_in`, SignIn).SetName(`sign_in`)
		g.Route(`GET,POST`, `/sign_out`, SignOut).SetName(`sign_out`)
		g.Route(`POST`, `/customer_info`, CustomerInfo).SetName(`customer_info`)
		g.Route(`POST`, `/qrcode/sign_in`, qrcodeSignIn).SetName(`qrcode_sign_in`)
		g.Route(`GET,POST`, `/forgot`, Forgot).SetName(`forgot`)
		g.Route(`GET,POST`, `/verification/callback/:provider/:recid/:timestamp/:token`, Verification)
		g.Route(`GET,POST`, `/custom/:page`, Custom).SetName(`custom`)
		g.Route(`GET,POST`, `/error_code`, ErrorCode).SetName(`error_code`)
		g.Route(`GET,POST`, `/secure_key`, SecureKey).SetName(`secure_key`)
		g.Route(`GET,POST`, `/advert/:idents`, Advert).SetName(`advert`)
	})
}
