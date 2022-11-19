package routepage

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/admpub/log"
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
		Register(v, router, funcMap)
	}
	return err
}

func Register(v *dbschema.OfficialCommonRoutePage, router *echo.Echo, funcMap map[string]interface{}) echo.IRouter {
	handler, err := buildHandler(v, funcMap)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info(`Register custom routing page: `, v.Name, ` (`, v.Method, ` `, v.Route, `)`)
	return router.Route(v.Method, v.Route, handler)
}

func buildHandler(v *dbschema.OfficialCommonRoutePage, funcMap map[string]interface{}) (func(ctx echo.Context) error, error) {
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
		}
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
