package user

import (
	"strings"

	"github.com/admpub/sessions"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

type QRSignIn struct {
	SessionID string
	sessionguard.Environment
}

func qrcodeScan(ctx echo.Context) error {
	if ctx.IsPost() {
		encrypted := ctx.Form(`data`)
		encrypted = strings.TrimSpace(encrypted)
		if len(encrypted) == 0 {
			return ctx.NewError(code.InvalidParameter, `参数“%s”不能为空`, `data`).SetZone(`data`)
		}
		encrypted = com.URLSafeBase64(encrypted, false)
		plaintext := config.FromFile().Decode256(encrypted)
		if len(plaintext) == 0 {
			return ctx.NewError(code.InvalidParameter, `解密失败`).SetZone(`data`)
		}
		signInData := &QRSignIn{}
		err := com.JSONDecodeString(plaintext, &signInData)
		if err != nil {
			return err
		}
		sessStore := ss.Get(ctx.SessionOptions().Engine)
		if sessStore == nil {
			return ctx.NewError(code.Unsupported, `不支持session存储引擎: %s`, ctx.SessionOptions().Engine)
		}

		// sign-in
		customer := sessdata.Customer(ctx)
		customerM := modelCustomer.NewCustomer(ctx).DisabledSession(true)
		err = customerM.Get(nil, `id`, customer.Id)
		if err != nil {
			if err == db.ErrNoMoreRows {
				return ctx.NewError(code.UserNotFound, `用户不存在`)
			}
			return err
		}
		co := modelCustomer.NewCustomerOptions(nil)
		co.Name = customer.Name
		co.SignInType = `qrcode`
		err = customerM.FireSignInSuccess(co, co.SignInType)
		if err != nil {
			return err
		}
		customerCopy := customerM.ClearPasswordData()

		tmpCookieName := ctx.SessionOptions().Name + `TMP`
		session := sessions.NewSession(sessStore, tmpCookieName)
		session.IsNew = false
		session.ID = signInData.SessionID
		session.Values[`customer_env`] = &signInData.Environment
		session.Values[`customer`] = &customerCopy
		err = sessStore.Save(ctx, session)
		if err != nil {
			return err
		}
		ctx.SetCookie(tmpCookieName, ``, -1)
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`操作成功`)))
	}
	return ctx.Render(`user/qrcode/scan`, nil)
}
