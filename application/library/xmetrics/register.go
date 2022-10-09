package xmetrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/webx-top/echo"
)

func Register(r echo.RouteRegister) {
	r.Use(RequestMiddleware())
	r.Get("/metrics", echo.WrapHandler(promhttp.Handler()), ResponseMiddleware())
}
