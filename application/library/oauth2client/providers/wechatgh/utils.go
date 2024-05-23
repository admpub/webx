package wechatgh

import (
	"sort"
	"strings"

	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func MakeSignature(timestamp, nonce, token string) string {
	arr := []string{timestamp, nonce, token}
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(arr); i++ {
		b.WriteString(arr[i])
	}

	return com.Sha1(b.String())
}

func CheckSignature(signature, timestamp, nonce, token string) bool {
	return MakeSignature(timestamp, nonce, token) == signature
}

func GetServer(ctx echo.Context, cfg *config.Config) *server.Server {
	officialAccount := GetAccount(cfg)
	// 传入request和responseWriter
	server := officialAccount.GetServer(ctx.Request().StdRequest(), ctx.Response().StdResponseWriter())
	return server
}
