package index

import (
	"github.com/admpub/nging/v5/application/handler/manager"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/webx-top/echo"
)

func init() {
	frontend.Register(func(g echo.RouteRegister) {
		g.Route(`GET,HEAD`, uploadLibrary.UploadURLPath+`:subdir/*`, manager.File)
		g.Route(`GET,POST`, `/`, Index)
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/search`, Search)
		g.Route(`GET,POST`, `/sign_up`, SignUp)
		g.Route(`GET,POST`, `/sign_in`, SignIn)
		g.Route(`GET,POST`, `/sign_out`, SignOut)
		g.Route(`POST`, `/customer_info`, CustomerInfo)
		g.Route(`GET,POST`, `/forgot`, Forgot)
		g.Route(`GET,POST`, `/verification/callback/:provider/:recid/:timestamp/:token`, Verification)
		g.Route(`GET,POST`, `/custom/:page`, Custom)
		g.Route(`GET,POST`, `/error_code`, ErrorCode)
		g.Route(`GET,POST`, `/secure_key`, SecureKey)
		g.Route(`GET,POST`, `/advert/:idents`, Advert)
	})
}
