package middleware

import (
	"net"

	"github.com/admpub/ipfilter"
	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config"
	cfgIPFilter "github.com/admpub/webx/application/library/ipfilter"
	caddyPluginIPfilter "github.com/caddy-plugins/ipfilter"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	mwIPFilter "github.com/webx-top/echo/middleware/ipfilter"
)

func init() {
	caddyPluginIPfilter.GetCountryCode = func(c caddyPluginIPfilter.IPFConfig, clientIP net.IP) (string, error) {
		if c.DBHandler == nil {
			return ipfilter.NetIPToCountry(clientIP), nil
		}
		// do the lookup.
		var result caddyPluginIPfilter.OnlyCountry
		err := c.DBHandler.Lookup(clientIP, &result)

		// get only the ISOCode out of the lookup results.
		return result.Country.ISOCode, err
	}

}

func getIPFilterConfig() (*cfgIPFilter.Options, bool) {
	opts, ok := config.Setting(`base`).Get(`ipFilter`).(*cfgIPFilter.Options)
	return opts, ok
}

var defaultIPFilterOptions = ipfilter.Options{
	AllowedIPs:       []string{},
	BlockedIPs:       []string{},
	AllowedCountries: []string{`CN`, `HK`, `TW`, `JP`, `RU`},
	BlockedCountries: []string{},
	Logger:           log.Writer(log.LevelInfo).(*log.LoggerWriter),
	TrustProxy:       true,
}

func IPFilter() echo.MiddlewareFuncd {
	ipfilterOptions := defaultIPFilterOptions
	if opts, ok := getIPFilterConfig(); ok {
		opts.Apply(&ipfilterOptions)
	}
	return mwIPFilter.IPFilter(mwIPFilter.Config{
		Skipper: func(c echo.Context) bool {
			opts, ok := getIPFilterConfig()
			if !ok || !opts.On || com.InSlice(c.RealIP(), cfgIPFilter.LocalIPs) {
				return true
			}
			passToken := c.Query(`passToken`)
			var fromCookie bool
			if len(passToken) == 0 {
				passToken = c.Cookie().Get(`passToken`)
				if len(passToken) == 0 {
					return false
				}
				fromCookie = true
			}
			if opts.PassToken != passToken {
				return false
			}
			if !fromCookie {
				c.Cookie().Set(`passToken`, passToken)
			}
			return true
		},
		Options: ipfilterOptions,
	})
}

func OnlyLocal() echo.MiddlewareFuncd {
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !com.IsLocalhost(c.Domain()) || !com.IsLocalhost(c.Request().RemoteAddress()) {
				return echo.ErrForbidden
			}
			return next.Handle(c)
		}
	}
}
