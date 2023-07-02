package cache

import (
	"context"

	"github.com/admpub/cache"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/db/lib/factory"
)

var cacheRootContext context.Context
var cacheRootContextCancel context.CancelFunc

func CloseCache() {
	if cacheRootContextCancel != nil {
		cacheRootContextCancel()
	}
}

func connection(diffs config.Diffs) error {
	cfg := common.Setting(`cache`)
	defaultCfg, ok := cfg.Get(`default`).(*cache.Options)
	if ok {
		CacheNew(cacheRootContext, *defaultCfg, `default`)
	}
	fallbackCfg, ok := cfg.Get(`fallback`).(*cache.Options)
	if ok {
		CacheNew(cacheRootContext, *fallbackCfg, `fallback`)
	}
	return nil
}

func init() {
	cacheRootContext, cacheRootContextCancel = context.WithCancel(context.Background())
	config.OnGroupSetSettings(`cache`, connection)
	factory.SetCacher(DBCacher)
}
