package xhtml

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/webx/application/library/cache"
	"github.com/admpub/webx/application/library/xcommon"
	"github.com/admpub/webx/application/middleware/sessdata"
	"github.com/admpub/webx/application/registry/route"
	"github.com/webx-top/com"
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
	body := rec.Body.String()
	err := cache.Put(context.Background(), saveAs, body+`<!-- Generated at `+time.Now().Format(time.DateTime)+` -->`, 0)
	if err != nil {
		log.Error(err)
	} else {
		err = cache.Put(context.Background(), saveAs+`.hash`, com.Md5(body), 0)
		if err != nil {
			log.Error(err)
		}
	}
	return err
}

func IsCached(ctx echo.Context, cacheKey string, urlWithQueryString ...bool) (bool, error) {
	if defaults.IsMockContext(ctx) {
		return false, nil
	}

	if customer := sessdata.Customer(ctx); customer != nil && customer.Uid > 0 {
		nocache := ctx.Formx(`nocache`).Int()
		switch nocache {
		case 1: // 禁用缓存
			return false, nil
		case 2: // 强制缓存新数据
			fallthrough
		case 3:
			reqURL := ctx.Request().URL().Path()
			if len(urlWithQueryString) > 0 && urlWithQueryString[0] {
				if query := ctx.Request().URL().RawQuery(); len(query) > 0 {
					reqURL += `?` + query
				}
			}
			if err := Make(http.MethodGet, reqURL, cacheKey); err != nil {
				return true, err
			}
		}
	}
	var cachedETag sql.NullString
	getHash := func() string {
		if !cachedETag.Valid {
			cache.Get(context.Background(), cacheKey+`.hash`, &cachedETag.String)
			cachedETag.Valid = true
		}
		return cachedETag.String
	}
	getContent := func() (string, error) {
		var cachedHTML string
		err := cache.Get(context.Background(), cacheKey, &cachedHTML)
		return cachedHTML, err
	}
	err := ETagCallback(ctx, getHash, getContent)
	return err == nil, err
}

func Remove(cacheKey string) error {
	cache.Delete(context.Background(), cacheKey+`.hash`)
	return cache.Delete(context.Background(), cacheKey)
}

func ETagCallback(ctx echo.Context, contentEtag func() string, callback func() (string, error), weak ...bool) error {
	var _weak bool
	if len(weak) > 0 {
		_weak = weak[0]
	}
	if reqETag := ctx.Header(`If-None-Match`); len(reqETag) > 0 {
		if _weak {
			reqETag = strings.TrimPrefix(reqETag, `W/`)
		}
		if reqETag == `"`+contentEtag()+`"` {
			return ctx.NotModified()
		}
	}
	cachedHTML, err := callback()
	if err != nil {
		return err
	}
	if len(contentEtag()) > 0 {
		eTag := `"` + contentEtag() + `"`
		if _weak { // 弱ETag是指在资源内容发生变化时，ETag值不一定会随之改变的标识符
			eTag = `W/` + eTag
		}
		ctx.Response().Header().Set(`ETag`, eTag)
	}
	return ctx.HTML(cachedHTML)
}
