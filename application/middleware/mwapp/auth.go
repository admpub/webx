package mwapp

import (
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/webx/application/library/xcode"
)

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{}
}

type AuthConfig struct {
	// Header 字段名
	HeaderAppIDKey string
	HeaderSignKey  string
	HeaderTimeKey  string
	// Form 字段名
	FormAppIDKey string
	FormSignKey  string
	FormTimeKey  string
	// 有效期
	LifeSeconds int64

	secretGetter func(ctx echo.Context, appID string) (string, error)
	signMaker    func(data url.Values, secret string) string
}

func (a *AuthConfig) setFormKV(ctx echo.Context, key string, value string) {
	ctx.Request().Form().Set(key, value)
}

func (a *AuthConfig) Prepare(ctx echo.Context, mustWithSign bool) (appID string, sign string, err error) {
	appID = ctx.Formx(a.FormAppIDKey).String()
	if len(appID) == 0 {
		appID = ctx.Header(a.HeaderAppIDKey)
		appID = strings.TrimSpace(appID)
		if len(appID) == 0 {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`无效参数值 %s: %s`, a.FormAppIDKey, appID)).SetZone(a.FormAppIDKey)
			return
		}
		a.setFormKV(ctx, a.FormAppIDKey, appID)
	}
	sign = ctx.Formx(a.FormSignKey).String()
	if len(sign) == 0 {
		sign = ctx.Header(a.HeaderSignKey)
		sign = strings.TrimSpace(sign)
		if len(sign) == 0 && mustWithSign {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`无效参数值 %s: %s`, a.FormSignKey, sign)).SetZone(a.FormSignKey)
			return
		}
		//a.setFormKV(ctx, a.FormSignKey, sign)
	}
	timestamp := ctx.Formx(a.FormTimeKey).Int64()
	if timestamp <= 0 {
		ts := ctx.Header(a.HeaderTimeKey)
		timestamp = param.AsInt64(ts)
		if timestamp <= 0 {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`无效参数值 %s: %d`, a.FormTimeKey, timestamp)).SetZone(a.FormTimeKey)
			return
		}
		a.setFormKV(ctx, a.FormTimeKey, ts)
	}
	if len(sign) > 0 && a.LifeSeconds > 0 {
		if time.Now().Unix()-timestamp > a.LifeSeconds {
			err = ctx.NewError(xcode.SignatureHasExpired, ctx.T(`页面已经失效，请返回重新提交`)).SetZone(a.FormTimeKey)
		}
	}
	return
}

func (a *AuthConfig) Verify(ctx echo.Context) error {
	appID, sign, err := a.Prepare(ctx, true)
	if err != nil {
		return err
	}
	gensign, err := a.SignRequest(ctx, appID)
	if err != nil {
		return err
	}
	//echo.Dump(echo.H{`sign`: sign, `make`: gensign})
	if sign != gensign {
		return ctx.NewError(stdCode.InvalidSignature, ctx.T(`签名无效`)).SetZone(a.FormSignKey)
	}
	return err
}

func (a *AuthConfig) SetSecretGetter(getter func(ctx echo.Context, appID string) (string, error)) *AuthConfig {
	a.secretGetter = getter
	return a
}

func (a *AuthConfig) SecretGetter() func(ctx echo.Context, appID string) (string, error) {
	return a.secretGetter
}

func (a *AuthConfig) SetSignMaker(signMaker func(data url.Values, secret string) string) *AuthConfig {
	a.signMaker = signMaker
	return a
}

func (a *AuthConfig) SignMaker() func(data url.Values, secret string) string {
	return a.signMaker
}

var (
	DefaultAuthAppConfig = AuthConfig{
		HeaderAppIDKey: `X-App-ID`,
		HeaderSignKey:  `X-App-Sign`,
		HeaderTimeKey:  `X-App-Timestamp`,
		FormAppIDKey:   `appID`,
		FormSignKey:    `sign`,
		FormTimeKey:    `timestamp`,
		LifeSeconds:    3600,
		secretGetter: func(ctx echo.Context, appID string) (secret string, err error) {
			return "", ctx.NewError(stdCode.Unsupported, ctx.T(`未设置App密钥获取方式！不支持获取App密钥`)).SetZone(`secretGetter`)
		},
		signMaker: func(data url.Values, secret string) string {
			return com.Sha256(data.Encode() + `&secret=` + secret)
		},
	}
)

func (a *AuthConfig) SetDefaults() {
	if len(a.HeaderAppIDKey) == 0 {
		a.HeaderAppIDKey = DefaultAuthAppConfig.HeaderAppIDKey
	}
	if len(a.HeaderSignKey) == 0 {
		a.HeaderSignKey = DefaultAuthAppConfig.HeaderSignKey
	}
	if len(a.HeaderTimeKey) == 0 {
		a.HeaderTimeKey = DefaultAuthAppConfig.HeaderTimeKey
	}
	if len(a.FormAppIDKey) == 0 {
		a.FormAppIDKey = DefaultAuthAppConfig.FormAppIDKey
	}
	if len(a.FormSignKey) == 0 {
		a.FormSignKey = DefaultAuthAppConfig.FormSignKey
	}
	if len(a.FormTimeKey) == 0 {
		a.FormTimeKey = DefaultAuthAppConfig.FormTimeKey
	}
	if a.secretGetter == nil {
		a.secretGetter = DefaultAuthAppConfig.secretGetter
	}
	if a.signMaker == nil {
		a.signMaker = DefaultAuthAppConfig.signMaker
	}
}

func Auth(_cfg AuthConfig) echo.MiddlewareFuncd {
	cfg := &_cfg
	cfg.SetDefaults()
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := cfg.Verify(c); err != nil {
				return err
			}
			return next.Handle(c)
		}
	}
}
