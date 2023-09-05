package official

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware"

	"github.com/admpub/webx/application/dbschema"
)

func NewRouteRewrite(ctx echo.Context) *RouteRewrite {
	m := &RouteRewrite{
		OfficialCommonRouteRewrite: dbschema.NewOfficialCommonRouteRewrite(ctx),
	}
	return m
}

type RouteRewrite struct {
	*dbschema.OfficialCommonRouteRewrite
}

func (f *RouteRewrite) check() error {
	ctx := f.Context()
	f.Name = strings.TrimSpace(f.Name)
	f.Route = strings.TrimSpace(f.Route)
	f.RewriteTo = strings.TrimSpace(f.RewriteTo)
	if len(f.Route) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入路由网址`).SetZone(`route`)
	}
	if len(f.RewriteTo) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入路由网址`).SetZone(`rewriteTo`)
	}
	if !strings.HasPrefix(f.Route, `/`) {
		f.Route = "/" + f.Route
	}
	if !strings.HasPrefix(f.RewriteTo, `/`) {
		f.RewriteTo = "/" + f.RewriteTo
	}
	err := middleware.ValidateRewriteRule(f.Route, f.RewriteTo)
	return err
}

func (f *RouteRewrite) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCommonRouteRewrite.Insert()
}

func (f *RouteRewrite) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCommonRouteRewrite.Update(mw, args...)
}
