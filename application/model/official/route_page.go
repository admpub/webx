package official

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	"github.com/admpub/log"
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
		return f.Context().E(`请输入名称`)
	}
	f.Route = strings.TrimSpace(f.Route)
	if len(f.Route) == 0 {
		return f.Context().E(`请输入路由网址`)
	}
	if !strings.HasPrefix(f.Route, `/`) {
		f.Route = "/" + f.Route
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

func (f *RoutePage) Register(v *dbschema.OfficialCommonRoutePage, router *echo.Echo) echo.IRouter {
	handler, err := f.buildHandler(v)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info(`Register custom routing page: `, v.Name, ` (`, v.Method, ` `, v.Route, `)`)
	return router.Route(v.Method, v.Route, handler)
}

func (f *RoutePage) buildHandler(v *dbschema.OfficialCommonRoutePage) (func(ctx echo.Context) error, error) {
	var (
		t      *template.Template
		parsed bool
		data   = echo.H{}
		parse  func(ctx echo.Context) error
	)
	if len(v.PageVars) > 0 {
		err := com.JSONDecode(com.Str2bytes(v.PageVars), &data)
		if err != nil {
			v.PageContent = `[buildHandler]parse page vars: ` + err.Error()
		} else {
			name := com.Md5(v.PageContent)
			t = template.New(name)
			parse = func(ctx echo.Context) (err error) {
				defer func() {
					if e := recover(); e != nil {
						if !strings.HasSuffix(fmt.Sprint(e), `cannot Parse after Execute`) {
							panic(e)
						}
						err = nil
					}
				}()
				_, err = t.Parse(v.PageContent)
				return
			}
		}
	}
	return func(ctx echo.Context) error {
		var content string
		if t != nil {
			t.Funcs(ctx.Funcs())
			if !parsed {
				if err := parse(ctx); err != nil {
					return errors.WithMessage(err, `[buildHandler]parse page content`)
				}
				parsed = true
			}
			buf := bytes.NewBuffer(nil)
			err := t.Execute(buf, data)
			if err != nil {
				content = err.Error()
			} else {
				content = buf.String()
			}
		} else {
			content = v.PageContent
		}
		switch v.PageType {
		case `redirect`:
			return ctx.Redirect(content)
		case `html`:
			return ctx.HTML(content)
		case `xml`:
			return ctx.XMLBlob(com.Str2bytes(content))
		case `json`:
			return ctx.JSONBlob(com.Str2bytes(content))
		default: //text
			return ctx.String(content)
		}
	}, nil
}

func (f *RoutePage) Apply(router *echo.Echo) error {
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, `N`)
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`id`)
	}, 0, -1, cond.And())
	for _, v := range f.Objects() {
		f.Register(v, router)
	}
	return err
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
