package wechatgh

import (
	"os"
	"testing"

	"github.com/admpub/godotenv"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
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
	err = MessageSystem(ctx, cfg, defaultHandler)
	assert.NoError(t, err)
}
