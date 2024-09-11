package cache

import (
	"context"
	"path/filepath"

	"github.com/admpub/cache"
	_ "github.com/admpub/cache/memcache" // memcache
	_ "github.com/admpub/cache/redis5"   // redis
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"

	dbschemaDBMgr "github.com/nging-plugins/dbmanager/application/dbschema"
)

var (
	defaultCacheOptions = &cache.Options{
		Adapter:  `memory`,
		Interval: 300,
	}
	defaultCacheInstance   cache.Cache
	cacheConfigParsers     = map[string]func(cache.Options) (cache.Options, error){}
	instanceCachePrefix    = `Cache:`
	defaultConnectionName  = `default`
	fallbackConnectionName = `fallback`
)

func AddOptionParser(adapter string, parser func(cache.Options) (cache.Options, error)) {
	cacheConfigParsers[adapter] = parser
}

func Cache(ctx context.Context, args ...string) cache.Cache {
	if len(args) > 0 {
		defaultConnectionName = args[0]
		if len(args) > 1 {
			fallbackConnectionName = args[2]
		}
	}
	key := instanceCachePrefix + defaultConnectionName
	c, ok := echo.Get(key).(cache.Cache)
	if ok {
		return c
	}
	logPrefix := color.GreenString(`[cache]`)
	log.Debug(logPrefix, `[`+defaultConnectionName+`] 未找到已连接的实例`)
	if defaultConnectionName != fallbackConnectionName {
		key = instanceCachePrefix + fallbackConnectionName
		c, ok = echo.Get(key).(cache.Cache)
		if ok {
			return c
		}
		log.Debug(logPrefix, `[`+fallbackConnectionName+`] 未找到已连接的实例`)
	}
	if c == nil {
		log.Debug(logPrefix, `[`+defaultCacheOptions.Adapter+`] 使用默认实例`)
		if defaultCacheInstance == nil {
			if ctx == nil {
				ctx = context.Background()
			}
			var err error
			defaultCacheInstance, err = cache.Cacher(ctx, *defaultCacheOptions)
			if err != nil {
				log.Errorf(logPrefix, `[`+defaultCacheOptions.Adapter+`] 使用默认实例错误: %v`, err)
			}
		}
		c = defaultCacheInstance
	}
	return c
}

func CacheNew(ctx context.Context, opts cache.Options, keys ...string) error {
	var connectionName string
	if len(keys) > 0 {
		connectionName = keys[0]
	}
	if len(connectionName) == 0 {
		connectionName = opts.Adapter
	}
	key := instanceCachePrefix + connectionName
	logPrefix := color.GreenString(`[cache]`) + `[` + connectionName + `][` + opts.Adapter + `]`
	c, ok := echo.Get(key).(cache.Cache)
	if ok {
		log.Info(logPrefix, `断开连接`)
		echo.Delete(key)
		c.Close()
	}
	if len(opts.AdapterConfig) == 0 {
		return nil
	}
	echo.FireByNameWithMap(`webx.cache.connected.`+opts.Adapter+`.before`, echo.H{`cache`: nil, `options`: opts})
	log.Info(logPrefix, `开始连接`)
	switch opts.Adapter {
	case `file`:
		if !filepath.IsAbs(opts.AdapterConfig) {
			opts.AdapterConfig = filepath.Join(echo.Wd(), opts.AdapterConfig)
		}
	case `redis`:
		if IsDbAccount(opts.AdapterConfig) {
			m := dbschemaDBMgr.NewNgingDbAccount(common.NewMockContext())
			err := m.Get(nil, `id`, opts.AdapterConfig)
			if err == nil {
				if len(m.Name) == 0 {
					m.Name = `0`
				}
				opts.AdapterConfig = `network=tcp,addr=` + m.Host + `,password=` + m.Password + `,db=` + m.Name + `,pool_size=100,idle_timeout=180,hset_name=Cache,prefix=cache:`
			}
		}
	}
	if cfgParser, ok := cacheConfigParsers[opts.Adapter]; ok {
		var err error
		opts, err = cfgParser(opts)
		if err != nil {
			log.Error(logPrefix, color.RedString(`配置解析失败:`+err.Error()))
			return err
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	val, err := cache.Cacher(ctx, opts)
	if err != nil {
		log.Error(logPrefix, color.RedString(`连接失败:`+err.Error()))
		return err
	}
	echo.Set(key, val)
	return echo.FireByNameWithMap(`webx.cache.connected.`+opts.Adapter+`.after`, echo.H{`cache`: val, `options`: opts})
}
