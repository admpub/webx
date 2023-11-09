package official

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
)

var RoutePageTypes = echo.NewKVData()

func init() {
	RoutePageTypes.Add(`html`, `HTML`)
	RoutePageTypes.Add(`json`, `JSON`)
	RoutePageTypes.Add(`text`, `纯文本`)
	RoutePageTypes.Add(`xml`, `XML`)
	RoutePageTypes.Add(`redirect`, `跳转网址`)
}

func NewRoutePage(ctx echo.Context) *RoutePage {
	m := &RoutePage{
		OfficialCommonRoutePage: dbschema.NewOfficialCommonRoutePage(ctx),
	}
	return m
}

type RoutePage struct {
	*dbschema.OfficialCommonRoutePage
}

func (f *RoutePage) check() error {
	f.Method = strings.TrimSpace(f.Method)
	f.Method = strings.Trim(f.Method, `,`)
	f.Name = strings.TrimSpace(f.Name)
	if len(f.Name) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入名称`).SetZone(`name`)
	}
	f.Route = strings.TrimSpace(f.Route)
	if len(f.Route) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入路由网址`).SetZone(`route`)
	}
	if !strings.HasPrefix(f.Route, `/`) {
		f.Route = "/" + f.Route
	}
	if f.PageType == `redirect` {
		f.TemplateEnabled = common.BoolN
	}
	return nil
}

func (f *RoutePage) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCommonRoutePage.Insert()
}

func (f *RoutePage) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCommonRoutePage.Update(mw, args...)
}

func (f *RoutePage) ListWithExtensionRoutes(ext string) ([]string, error) {
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, `N`)
	cond.AddKV(`route`, db.Like(`%`+ext))
	_, err := f.ListByOffset(nil, nil, 0, -1, cond.And())
	var routes []string
	for _, v := range f.Objects() {
		routes = append(routes, v.Route)
	}
	return routes, err
}
