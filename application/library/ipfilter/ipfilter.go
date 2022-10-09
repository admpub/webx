package ipfilter

import (
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/ipfilter"
)

func NewOptions() *Options {
	return &Options{}
}

type Options struct {
	On               bool     //是否启用
	PassToken        string   //通行口令
	BlockByDefault   bool     //默认封锁
	AllowedIPs       []string //白名单IP
	BlockedIPs       []string //黑名单IP
	AllowedCountries []string //白名单国家
	BlockedCountries []string //黑名单国家
	TrustProxy       bool     // 是否检查代理IP
}

func (o *Options) FromStore(r echo.H) *Options {
	o.On = r.Bool(`On`)
	o.PassToken = r.String(`PassToken`)
	o.BlockByDefault = r.Bool(`BlockByDefault`)
	o.TrustProxy = r.Bool(`TrustProxy`)
	allowedIPs := r.String(`AllowedIPs`)
	allowedIPs = strings.TrimSpace(allowedIPs)
	sep := "\n"
	if len(allowedIPs) > 0 {
		for _, v := range strings.Split(allowedIPs, sep) {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			o.AllowedIPs = append(o.AllowedIPs, v)
		}
	}
	blockedIPs := r.String(`BlockedIPs`)
	blockedIPs = strings.TrimSpace(blockedIPs)
	if len(blockedIPs) > 0 {
		for _, v := range strings.Split(blockedIPs, sep) {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			o.BlockedIPs = append(o.BlockedIPs, v)
		}
	}
	allowedCountries := r.String(`AllowedCountries`)
	allowedCountries = strings.TrimSpace(allowedCountries)
	if len(allowedCountries) > 0 {
		for _, country := range strings.Split(allowedCountries, sep) {
			country = strings.TrimSpace(country)
			if len(country) == 0 {
				continue
			}
			for _, v := range strings.Split(country, `,`) {
				v = strings.TrimSpace(v)
				if len(v) == 0 {
					continue
				}
				o.AllowedCountries = append(o.AllowedCountries, v)
			}
		}
	}
	blockedCountries := r.String(`BlockedCountries`)
	blockedCountries = strings.TrimSpace(blockedCountries)
	if len(blockedCountries) > 0 {
		for _, country := range strings.Split(blockedCountries, sep) {
			country = strings.TrimSpace(country)
			if len(country) == 0 {
				continue
			}
			for _, v := range strings.Split(country, `,`) {
				v = strings.TrimSpace(v)
				if len(v) == 0 {
					continue
				}
				o.BlockedCountries = append(o.BlockedCountries, v)
			}
		}
	}
	return o
}

func (o *Options) Apply(opts *ipfilter.Options) *Options {
	opts.AllowedIPs = o.AllowedIPs
	opts.BlockedIPs = o.BlockedIPs
	opts.AllowedCountries = o.AllowedCountries
	opts.BlockedCountries = o.BlockedCountries
	opts.BlockByDefault = o.BlockByDefault
	opts.TrustProxy = o.TrustProxy
	return o
}

var LocalIPs = []string{`127.0.0.1`, `::1`}
