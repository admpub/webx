package index

import (
	"io"
	"time"

	"github.com/coscms/sms/providers/twilio"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/model"
	uploadChecker "github.com/admpub/nging/v4/application/registry/upload/checker"
)

// Verification /verification/callback/:provider/:recid/:timestamp/:token
func Verification(c echo.Context) error {
	timestamp := c.Paramx(`timestamp`).Int64()
	if timestamp <= 0 {
		return c.E(`timestamp error`)
	}
	provider := c.Param(`provider`)     // 验证方式
	recid := c.Paramx(`recid`).Uint64() // 验证码记录ID
	token := c.Param(`token`)           // 签名
	if uploadChecker.Token(provider, recid, timestamp) != token {
		return c.E(`token error`)
	}
	if time.Now().Unix()-timestamp > 86400 {
		return c.E(`expired`)
	}
	m := model.NewVerification(c)
	cond := db.Cond{`id`: recid}
	err := m.Get(nil, cond)
	if err != nil {
		return err
	}
	switch provider {
	case `twilio`:
		body := c.Request().Body()
		b, e := io.ReadAll(body)
		if e != nil {
			return e
		}
		body.Close()
		status := `queued`
		result := &twilio.Notify{}
		err = com.JSONDecode(b, result)
		if err != nil {
			return err
		}
		if result.IsSuccess() {
			status = `success`
		} else if result.IsFailure() {
			status = `failure`
		}
		err = m.UpdateFields(nil, echo.H{
			`send_result`: string(b),
			`send_status`: status,
		}, cond)
	default:
	}
	return c.String(`OK`)
}
