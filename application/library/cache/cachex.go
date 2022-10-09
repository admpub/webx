package cache

import (
	"github.com/admpub/cache/x"
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

type (
	QueryFunc = x.QueryFunc
)

const (
	Disabled = x.Disabled
	Fresh    = x.Fresh
)

var (
	TTL               = x.TTL
	TTLFresh          = TTL(x.Fresh)
	TTLDisabled       = TTL(x.Disabled)
	Query             = x.Query
	DisableCacheUsage = x.DisableCacheUsage
	UseFreshData      = x.UseFreshData
)

func AdminRefreshable(ctx echo.Context, customer *dbschema.OfficialCustomer, ttl x.GetOption) x.GetOption {
	if customer == nil {
		return ttl
	}
	if customer.Uid <= 0 {
		return ttl
	}
	nocache := ctx.Formx(`nocache`).Bool()
	if !nocache {
		return ttl
	}
	return TTLFresh
}

func GenOptions(ctx echo.Context, cacheSeconds int64) []x.GetOption {
	nocache := ctx.Formx(`nocache`).Int()
	opts := []x.GetOption{TTL(cacheSeconds)}
	switch nocache {
	case 1:
		opts = append(opts, DisableCacheUsage(true)) // 禁用缓存
	case 2:
		opts = append(opts, UseFreshData(true)) // 强制缓存新数据
	case 3:
		opts = append(opts, UseFreshData(true))
	}
	return opts
}

func xNew(query x.Querier, ttlSeconds int64, args ...string) *x.Cachex {
	c := Cache(args...)
	return x.New(c, query, ttlSeconds)
}

// XQuery 获取缓存，如果不存在则执行函数获取数据并缓存【自动避免缓存穿透】
func XQuery(key string, recv interface{}, query x.Querier, options ...x.GetOption) error {
	return xNew(query, 0).Get(key, recv, options...)
}

// XFunc 获取缓存，如果不存在则执行函数获取数据并缓存【自动避免缓存穿透】
func XFunc(key string, recv interface{}, fn func() error, options ...x.GetOption) error {
	return xNew(QueryFunc(fn), 0).Get(key, recv, options...)
}

// Delete 删除缓存
func Delete(key string, args ...string) error {
	return Cache(args...).Delete(key)
}

// Put 设置缓存
func Put(key string, val interface{}, timeout int64, args ...string) error {
	return Cache(args...).Put(key, val, timeout)
}

// Get 获取缓存
func Get(key string, recv interface{}, args ...string) error {
	return Cache(args...).Get(key, recv)
}

// Incr 自增
func Incr(key string, args ...string) error {
	return Cache(args...).Incr(key)
}

// Decr 自减
func Decr(key string, args ...string) error {
	return Cache(args...).Decr(key)
}

// IsExist 是否存在缓存
func IsExist(key string, args ...string) bool {
	return Cache(args...).IsExist(key)
}

// Flush 删除所有缓存数据
func Flush(args ...string) error {
	return Cache(args...).Flush()
}
