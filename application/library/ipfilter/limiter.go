package ipfilter

import (
	"crypto/tls"
	"time"

	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/ratelimiter"
	"github.com/webx-top/echo/param"
	"gopkg.in/redis.v5"

	dbschemaDBMgr "github.com/nging-plugins/dbmanager/application/dbschema"
)

func NewRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{}
}

type RateLimiterConfig struct {
	On bool
	// The max count in duration for no policy, default is 100.
	Max int

	// Count duration for no policy, default is 1 Minute (60s).
	Duration int64 //seconds
	//key prefix, default is "LIMIT:".
	Prefix string

	//If request gets a  internal limiter error, just skip the limiter and let it go to next middleware
	SkipInternalError bool

	RedisAddr     string
	RedisPassword string
	RedisDB       int
	DBAccountID   uint
}

func (o *RateLimiterConfig) FromStore(r echo.H) *RateLimiterConfig {
	o.On = r.Bool(`On`)
	o.Max = r.Int(`Max`)
	o.Duration = r.Int64(`Duration`)
	o.Prefix = r.String(`Prefix`)
	o.SkipInternalError = r.Bool(`SkipInternalError`)
	o.RedisAddr = r.String(`RedisAddr`)
	o.DBAccountID = r.Uint(`DBAccountID`)
	return o
}

func (o *RateLimiterConfig) Apply(opts *ratelimiter.RateLimiterConfig) *RateLimiterConfig {
	opts.Max = o.Max
	opts.Duration = time.Duration(o.Duration) * time.Second
	opts.Prefix = o.Prefix
	opts.SkipRateLimiterInternalError = o.SkipInternalError
	if o.DBAccountID > 0 {
		m := dbschemaDBMgr.NewNgingDbAccount(defaults.NewMockContext())
		err := m.Get(nil, `id`, o.DBAccountID)
		if err == nil {
			if len(m.Name) == 0 {
				m.Name = `0`
			}
			o.RedisAddr = m.Host
			o.RedisPassword = m.Password
			o.RedisDB = param.AsInt(m.Name)
		}
	}
	if len(o.RedisAddr) > 0 {
		client := NewRedisClient(
			RedisAddr(o.RedisAddr),
			RedisPassword(o.RedisPassword),
			RedisDB(o.RedisDB),
		)
		if err := client.Ping().Err(); err != nil {
			log.Error(color.RedString(`[rateLimiter]`), ` `, err.Error())
		} else {
			opts.Client = client
		}
	}
	return o
}

func RedisAddr(addr string) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.Addr = addr
	}
}

func RedisMaxRetries(maxRetries int) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.MaxRetries = maxRetries
	}
}

func RedisNetwork(network string) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.Network = network
	}
}

func RedisPassword(password string) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.Password = password
	}
}

func RedisDB(db int) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.DB = db
	}
}

func RedisPoolSize(poolSize int) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.PoolSize = poolSize
	}
}

func RedisDialTimeout(timeout time.Duration) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.DialTimeout = timeout
	}
}

func RedisReadTimeout(timeout time.Duration) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.ReadTimeout = timeout
	}
}

func RedisWriteTimeout(timeout time.Duration) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.WriteTimeout = timeout
	}
}

func RedisPoolTimeout(timeout time.Duration) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.PoolTimeout = timeout
	}
}

func RedisIdleTimeout(timeout time.Duration) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.IdleTimeout = timeout
	}
}

func RedisIdleCheckFrequency(timeout time.Duration) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.IdleCheckFrequency = timeout
	}
}

func RedisReadOnly(readOnly bool) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.ReadOnly = readOnly
	}
}

func RedisTLSConfig(config *tls.Config) func(*redis.Options) {
	return func(opts *redis.Options) {
		opts.TLSConfig = config
	}
}

func NewRedisClient(settings ...func(*redis.Options)) *RedisClient {
	options := redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	}
	for _, option := range settings {
		option(&options)
	}
	c := &RedisClient{
		Client: redis.NewClient(&options),
	}
	return c
}

// RedisClient Implements RedisClient for redis.Client
type RedisClient struct {
	*redis.Client
}

func (c *RedisClient) DeleteKey(key string) error {
	return c.Del(key).Err()
}

func (c *RedisClient) EvalulateSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return c.EvalSha(sha1, keys, args...).Result()
}

func (c *RedisClient) LuaScriptLoad(script string) (string, error) {
	return c.ScriptLoad(script).Result()
}
