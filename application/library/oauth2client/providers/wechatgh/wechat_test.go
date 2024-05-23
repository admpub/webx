package wechatgh

import (
	"bytes"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/admpub/godotenv"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

func TestXxx(t *testing.T) {
	memory := cache.NewMemory()
	err := godotenv.Load(`./.env`)
	assert.NoError(t, err)
	cfg := &offConfig.Config{
		AppID:          os.Getenv(`AppID`),
		AppSecret:      os.Getenv(`AppSecret`),
		Token:          os.Getenv(`Token`),
		EncodingAESKey: os.Getenv(`EncodingAESKey`),
		Cache:          memory,
	}
	ctx := defaults.NewMockContext()
	timestamp := param.AsString(time.Now().Unix())
	nonce := com.RandomAlphanumeric(12)
	signature := MakeSignature(timestamp, nonce, cfg.Token)
	values := url.Values{}
	values.Set("signature", signature)
	values.Set("timestamp", timestamp)
	values.Set("nonce", nonce)
	ctx.Request().URL().SetRawQuery(values.Encode())

	ctx.Request().StdRequest().URL.RawQuery = values.Encode()
	assert.Equal(t, ctx.Request().StdRequest().URL.RawQuery, values.Encode())

	query := ctx.Request().StdRequest().URL.Query()
	assert.Equal(t, timestamp, query.Get(`timestamp`))

	reader := bytes.NewReader([]byte(`<xml><ToUserName>BUser</ToUserName><FromUserName>AUser</FromUserName><CreateTime>` + timestamp + `</CreateTime><MsgType>text</MsgType><Content>TEST</Content></xml>`))
	ctx.Request().SetBody(reader)
	err = MessageSystem(ctx, cfg, defaultHandler)
	assert.NoError(t, err)

	b := ctx.Response().Body()
	res := string(b)
	assert.Equal(t, res, `<xml><ToUserName><![CDATA[AUser]]></ToUserName><FromUserName><![CDATA[BUser]]></FromUserName><CreateTime>`+timestamp+`</CreateTime><MsgType>text</MsgType><Content><![CDATA[TEST]]></Content></xml>`)
}
