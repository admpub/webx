package top

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// StrengthenSafeSecret 强化的安全密钥(包含网络环境)
func StrengthenSafeSecret(ctx echo.Context, secret string) string {
	return secret + `|` + com.Md5(ctx.Request().UserAgent()) + `|` + ctx.RealIP()
}

// SecretSafeBuilder 密钥强化
func SecretSafeBuilder(ctx echo.Context, secret string) string {
	return secret + `|` + com.Md5(ctx.Request().UserAgent())
}
