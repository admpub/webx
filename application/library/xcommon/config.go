package xcommon

import (
	"strings"

	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/admpub/webx/application/library/frontend"
	"github.com/webx-top/echo"
)

var (
	BackendURL = common.BackendURL
)

func SiteURL(ctx echo.Context) string {
	return FrontendURL(ctx)
}

func FrontendURL(ctx echo.Context) string {
	frontendURL := common.Setting(`base`).String(`siteURL`)
	if len(frontendURL) == 0 {
		if ctx == nil {
			return frontendURL
		}
		frontendURL = ctx.Site()
	}
	frontendURL = strings.TrimSuffix(frontendURL, `/`) + strings.Trim(echo.String(`FrontendPrefix`), `/`)
	return frontendURL
}

func SiteName() string {
	return common.Setting(`base`).String(`siteName`)
}

func SiteSlogan() string {
	return common.Setting(`base`).String(`siteSlogan`)
}

func FrontendAssetsURL(ctx echo.Context) string {
	assetsURL := common.Setting(`base`, `assetsCDN`).String(`frontend`)
	if len(assetsURL) == 0 {
		assetsURL = FrontendURL(ctx) + frontend.AssetsURLPath
	}
	return assetsURL
}

func BackendAssetsURL(ctx echo.Context) string {
	assetsURL := common.Setting(`base`, `assetsCDN`).String(`backend`)
	if len(assetsURL) == 0 {
		assetsURL = BackendURL(ctx) + backend.AssetsURLPath
	}
	return assetsURL
}
