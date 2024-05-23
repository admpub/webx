package wechatgh

import (
	"context"
	"sync"

	"github.com/admpub/once"
	xcache "github.com/admpub/webx/application/library/cache"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
)

var wc *wechat.Wechat
var oc once.Once
var accounts = map[string]*officialaccount.OfficialAccount{}
var accMutex = &sync.RWMutex{}

func initWechat() {
	wc = wechat.NewWechat()
	ro := xcache.RedisOptions()
	if ro == nil {
		wc.SetCache(cache.NewMemory())
	} else {
		wc.SetCache(cache.NewRedis(context.Background(), &cache.RedisOpts{
			Host:     ro.Addr,
			Password: ro.Password,
			Database: ro.DB,
			//MaxIdle   :ro.MinIdleConns,
			IdleTimeout: int(ro.IdleTimeout.Seconds()),
		}))
	}
}

func GetWechat() *wechat.Wechat {
	oc.Do(initWechat)
	return wc
}

func GetAccount(cfg *config.Config) *officialaccount.OfficialAccount {
	accMutex.RLock()
	officialAccount, ok := accounts[cfg.AppID]
	accMutex.RUnlock()

	if ok {
		ctx := officialAccount.GetContext()
		if cfg.AppSecret == ctx.AppSecret &&
			cfg.EncodingAESKey == ctx.EncodingAESKey &&
			cfg.Token == ctx.Token &&
			cfg.Cache == ctx.Cache { // 配置完全相同的情况下
			return officialAccount
		}
	}

	wc := GetWechat()
	officialAccount = wc.GetOfficialAccount(cfg)
	accMutex.Lock()
	accounts[cfg.AppID] = officialAccount
	accMutex.Unlock()
	return officialAccount
}
