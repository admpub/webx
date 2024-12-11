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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/session"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func qrcodeScan(ctx echo.Context) error {
	if ctx.IsPost() {
		encrypted := ctx.Form(`data`)
		encrypted = strings.TrimSpace(encrypted)
		if len(encrypted) == 0 {
			return ctx.NewError(code.InvalidParameter, `参数“%s”不能为空`, `data`).SetZone(`data`)
		}
		if !strings.HasPrefix(encrypted, `coscms:`) {
			return ctx.NewError(code.InvalidParameter, `非本系统可以识别的参数`).SetZone(`data`)
		}
		encrypted = strings.TrimPrefix(encrypted, `coscms:`)
		parts := strings.SplitN(encrypted, `:`, 2)
		if len(parts) != 2 {
			return ctx.NewError(code.InvalidParameter, `无效的参数`).SetZone(`data`)
		}
		action := parts[0]
		decoder, ok := qrcodeDecoders[action]
		if !ok {
			return ctx.NewError(code.Unsupported, `不支持的数据操作“%s”`, action).SetZone(`data`)
		}
		encrypted = parts[1]
		if err := decoder(ctx, encrypted); err != nil {
			return err
		}
		return ctx.JSON(ctx.Data())
	}
	return ctx.Render(`user/qrcode/scan`, nil)
}

var qrcodeDecoders = map[string]func(ctx echo.Context, encrypted string) error{
	`signin`: qrcodeSignIn,
}

func RegisterQRCodeDecoder(name string, decoder func(ctx echo.Context, encrypted string) error) {
	qrcodeDecoders[name] = decoder
}

func qrcodeSignIn(ctx echo.Context, encrypted string) error {
	caseName := config.FromFile().Extend.GetStore(`QRCodeSignIn`).String(`case`)
	cs := GetQRSignInCase(caseName)
	signInData, err := cs.Decode(ctx, encrypted)
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
	ip2region.ClearZero(&ipInfo)
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
	ctx.Data().SetInfo(ctx.T(`扫码登录操作成功`))
	return err
}
