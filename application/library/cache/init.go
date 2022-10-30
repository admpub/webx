package cache

import (
	"github.com/admpub/cache"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/db/lib/factory"
)

func connection(diffs config.Diffs) error {
	cfg := common.Setting(`cache`)
	defaultCfg, ok := cfg.Get(`default`).(*cache.Options)
	if ok {
		CacheNew(*defaultCfg, `default`)
	}
	fallbackCfg, ok := cfg.Get(`fallback`).(*cache.Options)
	if ok {
		CacheNew(*fallbackCfg, `fallback`)
	}
	return nil
}

func init() {
	config.OnGroupSetSettings(`cache`, connection)
	factory.SetCacher(DBCacher)
}
