package mwapp

import (
	"bytes"
	"io"
	"net/url"

	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
)

func (a *AuthConfig) SignRequest(ctx echo.Context, appID string) (sign string, err error) {
	if a.secretGetter == nil {
		err = ctx.NewError(stdCode.Failure, ctx.T(`不支持获取密钥`))
		return
	}
	var secret string
	secret, err = a.secretGetter(ctx, appID)
	if err != nil {
		return
	}
	var data url.Values = ctx.Forms()
	switch ctx.ResolveContentType() {
	case echo.MIMEApplicationJSON, echo.MIMEApplicationXML:
		body := ctx.Request().Body()
		var b []byte
		b, err = io.ReadAll(body)
		body.Close()
		if err != nil {
			return
		}
		ctx.Request().SetBody(io.NopCloser(bytes.NewBuffer(b)))
		data.Set(`data`, string(b))
	}
	data.Del(a.FormSignKey)
	sign = a.signMaker(data, secret)
	return
}
