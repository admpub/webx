package ratelimiter

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/webx-top/echo"
)

type (
	// RateLimiterConfig defines the config for RateLimiter middleware.
	RateLimiterConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper

		// The max count in duration for no policy, default is 100.
		Max int

		// Count duration for no policy, default is 1 Minute.
		Duration time.Duration

		// Prefix key prefix, default is "LIMIT:".
		Prefix string

		// Use a redis client for limiter, if omit, it will use a memory limiter.
		Client RedisClient

		//If request gets a  internal limiter error, just skip the limiter and let it go to next middleware
		SkipRateLimiterInternalError bool

		LimiterKeyGenerator func(c echo.Context) (limiterKey string, policy []int)
	}

	limiter struct {
		abstractLimiter
		prefix string
	}

	// Result of limiter.Get
	Result struct {
		Total     int           // It Equals Options.Max, or policy max
		Remaining int           // It will always >= -1
		Duration  time.Duration // It Equals Options.Duration, or policy duration
		Reset     time.Time     // The limit record reset time
		Until     time.Duration
	}

	abstractLimiter interface {
		getLimit(key string, policy ...int) ([]interface{}, error)
		removeLimit(key string) error
	}

	//RedisClient interface
	RedisClient interface {
		DeleteKey(string) error
		EvalulateSha(string, []string, ...interface{}) (interface{}, error)
		LuaScriptLoad(string) (string, error)
	}
)

var (
	// DefaultRateLimiterConfig is the default rate limit middleware config.
	DefaultRateLimiterConfig = RateLimiterConfig{
		Skipper:                      echo.DefaultSkipper,
		Max:                          100,
		Duration:                     time.Minute * 1,
		Prefix:                       "LIMIT",
		Client:                       nil,
		SkipRateLimiterInternalError: false,
		LimiterKeyGenerator: func(c echo.Context) (string, []int) {
			return c.RealIP() + `@` + c.Method(), nil
		},
	}
)

// RateLimiter returns a rate limit middleware.
func RateLimiter() echo.MiddlewareFunc {
	return RateLimiterWithConfig(DefaultRateLimiterConfig)
}

// RateLimiterWithConfig returns a RateLimiter middleware with config.
// See: `RateLimiter()`.
func RateLimiterWithConfig(config RateLimiterConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimiterConfig.Skipper
	}

	if len(config.Prefix) == 0 {
		config.Prefix = "LIMIT:"
	}
	if config.Max <= 0 {
		config.Max = 100
	}
	if config.Duration <= 0 {
		config.Duration = time.Minute * 1
	}
	if config.LimiterKeyGenerator == nil {
		config.LimiterKeyGenerator = DefaultRateLimiterConfig.LimiterKeyGenerator
	}

	var limiterImp *limiter
	//If config.Client omit, the limiter is a memory limiter
	if config.Client == nil {
		limiterImp = newMemoryLimiter(&config)
	} else {
		//setup redis client
		limiterImp = newRedisLimiter(&config)
	}

	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) {
				return h.Handle(c)
			}

			//policy := []int{10,1000}
			/*custom policy will configurable like
			[
				{"RealIP+Method+RequestURI","Max Value","Duration"},
				{"RealIP+Method+RequestURI","Max Value","Duration"}
			]
			*/
			id, policy := config.LimiterKeyGenerator(c)
			result, err := limiterImp.Get(id, policy...)

			if err != nil {
				if config.SkipRateLimiterInternalError {
					return h.Handle(c)
				}
				return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
			}

			response := c.Response()
			response.Header().Set("X-Ratelimit-Limit", strconv.FormatInt(int64(result.Total), 10))
			response.Header().Set("X-Ratelimit-Remaining", strconv.FormatInt(int64(result.Remaining), 10))
			response.Header().Set("X-Ratelimit-Reset", strconv.FormatInt(result.Reset.Unix(), 10))

			if result.Remaining <= 0 {
				until := result.Until
				after := int64(until) / 1e9
				response.Header().Set("Retry-After", strconv.FormatInt(after, 10))
				retryAfter := until.String()
				response.Header().Set("X-Retry-After", retryAfter)
				return echo.NewHTTPError(http.StatusTooManyRequests, fmt.Sprintf("Rate limit exceeded, retry in %s", retryAfter))
			}
			return h.Handle(c)
		})
	}
}

// get & remove

func (l *limiter) Get(id string, policy ...int) (Result, error) {
	var result Result
	key := l.prefix + id

	if odd := len(policy) % 2; odd == 1 {
		return result, errors.New("ratelimiter: must be paired values")
	}

	res, err := l.getLimit(key, policy...)
	if err != nil {
		return result, err
	}

	result = Result{}
	switch res[3].(type) {
	case time.Time: // result from memory limiter
		result.Remaining = res[0].(int)
		result.Total = res[1].(int)
		result.Duration = res[2].(time.Duration)
		result.Reset = res[3].(time.Time)
	default: // result from disteributed limiter
		result.Remaining = int(res[0].(int64))
		result.Total = int(res[1].(int64))
		result.Duration = time.Duration(res[2].(int64) * 1e6)

		timestamp := res[3].(int64)
		sec := timestamp / 1000
		result.Reset = time.Unix(sec, (timestamp-(sec*1000))*1e6)
	}
	if result.Remaining <= 0 {
		result.Until = time.Until(result.Reset)
	}
	return result, nil
}

// Remove remove limiter record for id
func (l *limiter) Remove(id string) error {
	return l.removeLimit(l.prefix + id)
}
