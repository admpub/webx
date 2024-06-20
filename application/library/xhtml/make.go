package xhtml

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/webx/application/library/cache"
	"github.com/admpub/webx/application/library/xcommon"
	"github.com/admpub/webx/application/registry/route"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	test "github.com/webx-top/echo/testing"
)

var ErrGenerateHTML = errors.New(`failed to generate html`)

func Make(method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	siteURL := xcommon.SiteURL(nil)
	if len(siteURL) == 0 {
		return fmt.Errorf(`%w: frontend URL cannot be empty`, ErrGenerateHTML)
	}
	rec := test.Request(method, siteURL+path, route.IRegister().Echo(), reqRewrite...)
	if rec.Code != http.StatusOK {
		return fmt.Errorf(`%w: [%d] %v`, ErrGenerateHTML, rec.Code, rec.Body.String())
	}
	err := cache.Put(context.Background(), saveAs, rec.Body.String()+`<!-- Generated at `+time.Now().Format(time.DateTime)+` -->`, 0)
	if err != nil {
		log.Error(err)
	}
	return err
}

func IsCached(ctx echo.Context, cacheKey string) (bool, error) {
	if defaults.IsMockContext(ctx) {
		return false, nil
	}
	var cachedHTML string
	err := cache.Get(context.Background(), cacheKey, &cachedHTML)
	if err == nil {
		return true, ctx.HTML(cachedHTML)
	}
	return false, err
}
