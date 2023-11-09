package routepage

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	xMW "github.com/admpub/webx/application/middleware"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func Apply(router *echo.Echo, funcMap map[string]interface{}) error {
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, `N`)
	f := dbschema.NewOfficialCommonRoutePage(defaults.NewMockContext())
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`id`)
	}, 0, -1, cond.And())
	for _, v := range f.Objects() {
		Register(*v, router, funcMap)
	}
	return err
}

func Register(v dbschema.OfficialCommonRoutePage, router *echo.Echo, funcMap map[string]interface{}) echo.IRouter {
	handler, err := buildHandler(v, funcMap)
	if err != nil {
		log.Errorf(`Failed to build custom routing page([%s]%s): %v`, v.Method, v.Route, err)
		return nil
	}
	log.Info(`Register custom routing page: `, v.Name, ` (`, v.Method, ` `, v.Route, `)`)
	return router.Route(v.Method, v.Route, handler)
}

func buildHandler(v dbschema.OfficialCommonRoutePage, funcMap map[string]interface{}) (func(ctx echo.Context) error, error) {
	var (
		t    *template.Template
		data = echo.H{}
		err  error
	)
	if len(v.PageVars) > 0 {
		err = com.JSONDecode(com.Str2bytes(v.PageVars), &data)
		if err != nil {
			v.PageContent = `[buildHandler]parse page vars: ` + err.Error()
		} else {
			name := com.Md5(v.PageContent)
			t = template.New(name)
			t.Funcs(funcMap)
			defer func() {
				if e := recover(); e != nil {
					err = fmt.Errorf(`%v`, e)
				}
			}()
			_, err = t.Parse(v.PageContent)
			if err != nil {
				err = echo.ParseTemplateError(err, v.PageContent)
			}
		}
	}

	tmplRender := func(ctx echo.Context, content interface{}) (bool, error) {
		if v.TemplateEnabled == common.BoolY && len(v.TemplateFile) > 0 {
			tmpl := v.TemplateFile
			tmpl = strings.TrimSuffix(tmpl, `.html`)
			if len(tmpl) > 0 {
				ctx.Set(`title`, v.Name)
				ctx.Set(`routePageContent`, content)
				return true, ctx.Render(path.Join(`route_page`, tmpl), nil)
			}
		}
		return false, nil
	}

	return func(ctx echo.Context) error {
		var content string
		if t != nil {
			buf := bytes.NewBuffer(nil)
			err := t.Execute(buf, xMW.DefaultRenderDataWrapper(ctx, data))
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
			if ok, err := tmplRender(ctx, template.HTML(content)); ok {
				return err
			}
			return ctx.HTML(content)
		case `xml`:
			if ok, err := tmplRender(ctx, template.HTML(content)); ok {
				return err
			}
			return ctx.XMLBlob(com.Str2bytes(content))
		case `json`:
			if ok, err := tmplRender(ctx, template.JS(content)); ok {
				return err
			}
			return ctx.JSONBlob(com.Str2bytes(content))
		default: //text
			if ok, err := tmplRender(ctx, content); ok {
				return err
			}
			return ctx.String(content)
		}
	}, err
}
