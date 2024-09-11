package cache

import (
	"github.com/admpub/cache"
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/coscms/webcore/library/config/extend"
)

func IsDbAccount(v string) bool {
	return com.StrIsNumeric(v)
}

var _ extend.Reloader = (*ReloadableOptions)(nil)

func NewReloadableOptions() *ReloadableOptions {
	return &ReloadableOptions{
		Options: &cache.Options{},
	}
}

type ReloadableOptions struct {
	*cache.Options
}

func (o *ReloadableOptions) Reload() error {
	err := CacheNew(cacheRootContext, *o.Options, `locker`)
	if err != nil {
		logPrefix := color.GreenString(`[cache]`) + `[locker][` + o.Adapter + `]`
		log.Error(logPrefix, err)
	} else {
		if o.Adapter == `redis` {
			resetRedsync()
			SetDefaultLockType(LockTypeRedis)
		} else {
			SetDefaultLockType(LockTypeMemory)
		}
	}
	return err
}

func (o *ReloadableOptions) IsValid() bool {
	return o != nil && o.Options != nil && len(o.Options.Adapter) > 0
}
