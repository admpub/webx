package frontend

import (
	"os"
	"path/filepath"
	"strings"

	godl "github.com/admpub/go-download/v2"
	"github.com/admpub/log"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var faviconCacheDir = filepath.Join(echo.Wd(), `data`, `cache`)
var faviconCacheFile = filepath.Join(faviconCacheDir, `favicon.ico`)

func init() {
	config.OnKeySetSettings(`base.siteFavicon`, func(diff config.Diff) error {
		siteFavicon := diff.String()
		n := len(siteFavicon)
		if !diff.IsDiff {
			if com.IsFile(faviconCacheFile) {
				if n == 0 {
					os.Remove(faviconCacheFile)
				}
				return nil
			}
		}
		if n > 7 {
			switch siteFavicon[0:7] {
			case `https:/`, `http://`:
				return faviconCache(siteFavicon)
			default:
				if siteFavicon[0:2] == `//` {
					return faviconCache(`http:` + siteFavicon)
				}
			}
		}
		if n > 0 {
			siteFavicon = strings.TrimPrefix(siteFavicon, `/`)
			if com.IsFile(siteFavicon) {
				com.MkdirAll(faviconCacheDir, os.ModePerm)
				return com.Copy(siteFavicon, faviconCacheFile)
			}
		}
		err := os.Remove(faviconCacheFile)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
		}
		return err
	})
}

func faviconCache(url string) (err error) {
	log.Debugf(`download favicon: %s => %s`, url, faviconCacheFile)
	_, err = godl.Download(url, faviconCacheFile, nil)
	return
}

func faviconHandler(c echo.Context) error {
	if com.IsFile(faviconCacheFile) {
		return c.File(faviconCacheFile)
	}
	return bootconfig.FaviconHandler(c)
}
