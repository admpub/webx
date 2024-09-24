package article

import (
	"net/http"
	"strings"
	"time"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/library/top"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func Redirect(ctx echo.Context) error {
	if ctx.IsPost() {
		return verifyRedirectURL(ctx)
	}
	customer := sessdata.Customer(ctx)
	if customer == nil {
		return nerrors.ErrUserNotLoggedIn
	}
	username := customer.Name
	url := ctx.Query(`url`)
	if len(url) == 0 {
		return ctx.NewError(code.InvalidParameter, `网址无效`).SetZone(`url`)
	}
	var expiry int64
	var err error
	url, expiry, err = top.ParseEncodedURL(url)
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
	parts := strings.SplitN(url, `#`, 2)
	parts[0] = com.WithURLParams(parts[0], `_token`, config.FromFile().Encode256(url+`||`+time.Now().Format(time.DateTime)+`||`+username), `customer`, username)
	url = strings.Join(parts, `#`)
	return ctx.Redirect(url, http.StatusFound)
}

// url, token, customer
func verifyRedirectURL(ctx echo.Context) error {
	url := ctx.Form(`url`)
	data := ctx.Data()
	if len(url) == 0 {
		data.SetError(ctx.NewError(code.InvalidParameter, `网址无效`).SetZone(`url`))
		return ctx.JSON(data)
	}
	username := ctx.Form(`customer`)
	token := ctx.Form(`token`)
	if len(token) == 0 {
		data.SetError(ctx.NewError(code.InvalidParameter, `令牌无效`).SetZone(`token`))
		return ctx.JSON(data)
	}
	content := config.FromFile().Decode256(token)
	if len(content) == 0 {
		data.SetError(ctx.NewError(code.DataFormatIncorrect, `令牌分析失败`).SetZone(`token`))
		return ctx.JSON(data)
	}
	parts := strings.SplitN(content, `||`, 3)
	if len(parts) != 3 {
		data.SetError(ctx.NewError(code.DataFormatIncorrect, `令牌分析后的格式不正确`).SetZone(`token`))
		return ctx.JSON(data)
	}
	if parts[0] != url {
		data.SetError(ctx.NewError(code.DataUnavailable, `令牌与网址不匹配`).SetZone(`url`))
		return ctx.JSON(data)
	}
	if len(parts[2]) > 0 && parts[2] != username {
		data.SetError(ctx.NewError(code.DataUnavailable, `令牌与客户不匹配`).SetZone(`customer`))
		return ctx.JSON(data)
	}
	t, err := time.Parse(time.DateTime, parts[1])
	if err != nil {
		data.SetError(err)
		return ctx.JSON(data)
	}
	if t.Before(time.Now().Add(time.Minute * -5)) {
		data.SetError(ctx.NewError(code.DataHasExpired, `令牌已经过期`).SetZone(`token`))
	} else {
		data.SetInfo(`OK`, code.Success.Int())
	}
	return ctx.JSON(data)
}
