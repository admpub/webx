package shorturl

import (
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	modelShorturl "github.com/admpub/webx/application/model/official/shorturl"
)

// Find 访问短网址
func Find(ctx echo.Context) error {
	shortID := ctx.Param(`shortId`)
	m := modelShorturl.NewShortURL(ctx)
	rawURL, err := m.Find(shortID)
	if err != nil {
		cause := errors.Cause(err)
		switch cause {
		case modelShorturl.ErrNotExistsURL:
			ctx.Set(`shortId`, shortID)
			return ctx.Render(`short_url/not_exists`, nil)
		case modelShorturl.ErrNeedURLPassword, modelShorturl.ErrWrongURLPassword:
			ctx.Set(`shortId`, shortID)
			if cause == modelShorturl.ErrNeedURLPassword && !ctx.IsPost() {
				err = nil
			}
			return ctx.Render(`short_url/verify`, err)
		case modelShorturl.ErrExpiredURL:
			ctx.Set(`shortId`, shortID)
			msg := err.Error()
			pos := strings.LastIndex(msg, `:`)
			if pos > 0 {
				err = errors.New(msg[0:pos])
			}
			ctx.Set(`title`, cause.Error())
			return ctx.Render(`short_url/expired`, nil)
		}
		return err
	}
	return ctx.Redirect(rawURL)
}
