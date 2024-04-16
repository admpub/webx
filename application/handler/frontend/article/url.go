package article

import (
	"net/http"
	"time"

	"github.com/admpub/webx/application/library/top"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func Redirect(ctx echo.Context) error {
	url := ctx.Query(`url`)
	url, expiry, err := top.ParseEncodedURL(url)
	if err != nil {
		return err
	}
	if len(url) == 0 {
		return ctx.NewError(code.Failure, `网址无效`)
	}
	if expiry == 0 {
		expiry = ctx.Queryx(`expiry`).Int64()
	}
	if expiry > 0 && time.Now().Unix() >= expiry {
		return ctx.NewError(code.DataHasExpired, `网址已经失效`)
	}
	return ctx.Redirect(url, http.StatusFound)
}
