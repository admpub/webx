package user

import (
	"strings"
	"time"

	"github.com/admpub/sessions"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/ip2region"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/session"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

type QRSignIn struct {
	SessionID     string `json:"sID"`
	SessionMaxAge int    `json:"sAge"`
	Expires       int64  `json:"dExp"` // 数据过期时间戳(秒)
	IPAddress     string `json:"ip"`
	UserAgent     string `json:"ua"`
	Platform      string `json:"pf"`
	Scense        string `json:"ss"`
	DeviceNo      string `json:"dn"`
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
		//echo.Dump(echo.H{`qrcodeScan`: signInData})
		if signInData.Expires < time.Now().Unix() {
			return ctx.NewError(code.DataHasExpired, `二维码已经过期`).SetZone(`data`)
		}
		sessStore := ss.Get(ctx.SessionOptions().Engine)
		if sessStore == nil {
			return ctx.NewError(code.Unsupported, `不支持 session 存储引擎: %s`, ctx.SessionOptions().Engine)
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
		co.ApplyOptions(
			modelCustomer.CustomerPlatform(signInData.Platform),
			modelCustomer.CustomerScense(signInData.Scense),
			modelCustomer.CustomerDeviceNo(signInData.DeviceNo),
			modelCustomer.CustomerMaxAgeSeconds(signInData.SessionMaxAge),
			modelCustomer.CustomerSessionID(signInData.SessionID),
			modelCustomer.CustomerIPAddress(signInData.IPAddress),
		)
		session.RememberMaxAge(ctx, signInData.SessionMaxAge)

		err = customerM.FireSignInSuccess(co, co.SignInType)
		if err != nil {
			return err
		}
		customerCopy := customerM.ClearPasswordData()

		tmpCookieName := ctx.SessionOptions().Name + `TMP`
		session := sessions.NewSession(sessStore, tmpCookieName)
		session.IsNew = false
		session.ID = signInData.SessionID

		// set session values
		ipInfo, _ := ip2region.IPInfo(signInData.IPAddress)
		sEnv := &sessionguard.Environment{
			UserAgent: signInData.UserAgent,
			Location:  ipInfo,
		}
		session.Values[`customer_env`] = sEnv
		session.Values[`customer`] = &customerCopy
		session.Values[`deviceInfo`] = &co.DeviceInfo

		err = sessStore.Save(ctx, session)
		if err != nil {
			return err
		}
		ctx.SetCookie(tmpCookieName, ``, -1)
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`扫码登录操作成功`)))
	}
	return ctx.Render(`user/qrcode/scan`, nil)
}
