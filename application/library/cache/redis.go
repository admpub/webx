package cache

import (
	"github.com/admpub/once"
	"github.com/webx-top/echo"
	"gopkg.in/redis.v5"
)

var (
	redisClient *redis.Client
	redisOnce   once.Once
)

func init() {
	echo.OnCallback(`webx.cache.connected.redis.after`, func(data echo.Event) error {
		resetRedisClient()
		return nil
	})
}

func resetRedisClient() {
	redisOnce.Reset()
}

func initRedisClient() {
	defer resetRedsync()
	rc, ok := Cache(cacheRootContext, `default`).Client().(*redis.Client)
	if ok {
		redisClient = rc
		return
	}
	rc, _ = Cache(cacheRootContext, `fallback`).Client().(*redis.Client)
	redisClient = rc
}

func onceInitRedisClient() {
	initRedisClient()
}

func RedisClient() *redis.Client {
	redisOnce.Do(onceInitRedisClient)
	return redisClient
}

func RedisOptions() *redis.Options {
	opt, ok := Cache(cacheRootContext, `default`).(redisOptions)
	if ok {
		return opt.Options()
	}
	opt, ok = Cache(cacheRootContext, `fallback`).(redisOptions)
	if ok {
		return opt.Options()
	}
	return nil
}

type redisOptions interface {
	Options() *redis.Options
}
