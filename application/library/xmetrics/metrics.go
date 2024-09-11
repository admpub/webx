package xmetrics

import (
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	echoPrometheus "github.com/webx-top/echo-prometheus"
	"github.com/webx-top/echo/middleware"
)

var (
	_ extend.Reloader = &Metrics{}
)

func New() *Metrics {
	return &Metrics{
		BasicAuth: &BasicAuth{},
	}
}

type Metrics struct {
	*BasicAuth
	Enable bool `json:"enable"`
}

func (m *Metrics) Reload() error {
	return nil
}

type BasicAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func RequestMiddleware() echo.Middleware {
	prometheusConfig := echoPrometheus.DefaultConfig
	prometheusConfig.Skipper = RequestSkipper()
	prometheusConfig.OnlyRoutePath = true
	return echoPrometheus.MetricsMiddlewareWithConfig(prometheusConfig)
}

func RequestSkipper() echo.Skipper {
	return func(c echo.Context) bool {
		mc, ok := config.FromFile().Extend.Get(Name).(*Metrics)
		if !ok || !mc.Enable {
			return true
		}
		return c.Request().URL().Path() == `/metrics`
	}
}

func ResponseMiddleware() echo.Middleware {
	return middleware.BasicAuth(func(user string, password string) bool {
		mc, ok := config.FromFile().Extend.Get(Name).(*Metrics)
		if !ok || mc.BasicAuth == nil {
			return false
		}
		if !mc.Enable || len(mc.User) == 0 || len(mc.Password) == 0 {
			return false
		}
		return mc.User == user && mc.Password == password
	}, func(c echo.Context) bool {
		return com.IsLocalhost(c.Domain()) && com.IsLocalhost(c.Request().RemoteAddress())
	})
}
