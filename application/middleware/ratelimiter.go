package middleware

import (
	"github.com/admpub/nging/v4/application/library/config"
	cfgIPFilter "github.com/admpub/webx/application/library/ipfilter"
	"github.com/webx-top/echo"
	mwRateLimiter "github.com/webx-top/echo/middleware/ratelimiter"
)

func getRateLimiterConfig() (*cfgIPFilter.RateLimiterConfig, bool) {
	opts, ok := config.Setting(`frequency`).Get(`rateLimiter`).(*cfgIPFilter.RateLimiterConfig)
	return opts, ok
}

func RateLimiter() echo.MiddlewareFunc {
	rateLimiterConfig := &mwRateLimiter.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			opts, ok := getRateLimiterConfig()
			if !ok || !opts.On {
				return true
			}
			return false
		},
	}
	if opts, ok := getRateLimiterConfig(); ok {
		opts.Apply(rateLimiterConfig)
	}
	return mwRateLimiter.RateLimiterWithConfig(*rateLimiterConfig)
}
