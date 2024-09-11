package middleware

import (
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/codec"
	"github.com/coscms/webcore/library/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	mwJWT "github.com/webx-top/echo/middleware/jwt"
	"github.com/webx-top/echo/param"
)

func defaultJWTSkipper(c echo.Context) bool {
	enabledJWT := c.Query(`client`) == `app` || c.Header(`X-Client`) == `app`
	c.Internal().Set(`enabledJWT`, enabledJWT)
	return !enabledJWT
}

func JWT(skippers ...func(echo.Context) bool) echo.MiddlewareFuncd {
	jwtConfig := &mwJWT.JWTConfig{
		Skipper:     defaultJWTSkipper,
		SigningKey:  []byte(config.FromFile().Cookie.HashKey),
		Claims:      &jwt.StandardClaims{},
		TokenLookup: "header:" + echo.HeaderAuthorization,
	}
	if len(skippers) > 0 && skippers[0] != nil {
		jwtConfig.Skipper = skippers[0]
	}
	jwtConfig.SetErrorHandler(func(ctx echo.Context, err error) {
		log.Debug(`JWT: `, err)
	}).SetFallbackExtractor(jwtQueryExtractor).SetTokenPreprocessor(jwtTokenPreprocessor)
	return mwJWT.JWTWithConfig(*jwtConfig)
}

// 从网址中传递时，需要验证ip和有效期
func jwtQueryExtractor(ctx echo.Context) (string, error) {
	token := ctx.Query(`jwt`)
	var err error
	if len(token) == 0 {
		return ``, mwJWT.ErrJWTMissing
	}
	expiry := ctx.Queryx(`exp`).Int64()
	if expiry < 1 {
		return ``, mwJWT.ErrJWTMissing
	}
	nowTS := time.Now().Unix()
	if nowTS > expiry {
		return ``, mwJWT.ErrJWTMissing
	}
	token, err = codec.DefaultSM2DecryptHex(token)
	if err != nil {
		return ``, ctx.NewError(code.InvalidParameter, `JWT解密失败: %v`, err)
	}
	parts := strings.SplitN(token, `|`, 3) // jwt|ip|ts
	if len(parts) != 3 {
		return ``, mwJWT.ErrJWTMissing
	}
	token = parts[0]
	ts := param.AsInt64(parts[2])
	if ts != expiry {
		return ``, mwJWT.ErrJWTMissing
	}
	ip := parts[1]
	if ip != ctx.RealIP() {
		return ``, mwJWT.ErrJWTMissing
	}
	return token, err
}

// 支持解密
func jwtTokenPreprocessor(ctx echo.Context, token string) (string, error) {
	if strings.HasPrefix(token, `{sm2}`) {
		token = strings.TrimPrefix(token, `{sm2}`)
		return codec.DefaultSM2DecryptHex(token)
	}
	return token, nil
}
