package {{.PkgName}}

import (
	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/webx/application/model/{{.Group}}"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func {{.H.Name}}Index(ctx echo.Context) error {
	m := {{.M.PkgName}}.New{{.M.Name}}(ctx)
	cond := db.NewCompounds()
	{{- if .M.NameField}}
	common.SelectPageCond(ctx, cond, `{{.M.IDColumn}}`, `{{.M.NameColumn}}%`)
	{{- end}}
	sorts := []interface{}{`-{{.M.IDColumn}}`}
	list, err := m.ListPage(cond, sorts...)
	ctx.Set(`listData`, list)
	ctx.Set(`title`, ctx.T(`%v列表`, {{.M.PkgName}}.ObjectName{{.M.Name}}))
	return ctx.Render(`{{.Group}}/{{.H.TmplName "index"}}`, handler.Err(ctx, err))
}

func {{.H.Name}}Add(ctx echo.Context) error {
	var err error
	m := {{.M.PkgName}}.New{{.M.Name}}(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.{{.M.SchemaName}})
		if err == nil {
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/{{.Group}}/{{.H.TmplName "index"}}`))
			}
		}
	} else {
		{{.M.IDField|LowerCaseFirst}} := ctx.Formx(`copy{{.M.IDField}}`).{{.M.IDFieldType|Title}}()
		if {{.M.IDField|LowerCaseFirst}} > 0 {
			err = m.Get(nil, `{{.M.IDColumn}}`, {{.M.IDField|LowerCaseFirst}})
			if err == nil {
				echo.StructToForm(ctx, m.{{.M.SchemaName}}, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`{{.M.IDColumn}}`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/{{.Group}}/{{.H.TmplName "index"}}`)
	ctx.Set(`title`, ctx.T(`添加%v`, {{.M.PkgName}}.ObjectName{{.M.Name}}))
	ctx.Set(`listPageTitle`, ctx.T(`%v列表`, {{.M.PkgName}}.ObjectName{{.M.Name}}))
	return ctx.Render(`{{.Group}}/{{.H.TmplName "edit"}}`, err)
}

func {{.H.Name}}Edit(ctx echo.Context) error {
	var err error
	{{.M.IDField|LowerCaseFirst}} := ctx.Formx(`{{.M.IDField|LowerCaseFirst}}`).{{.M.IDFieldType|Title}}()
	m := {{.M.PkgName}}.New{{.M.Name}}(ctx)
	err = m.Get(nil, db.Cond{`{{.M.IDColumn}}`: {{.M.IDField|LowerCaseFirst}}})
	if ctx.IsPost() {
		err = ctx.MustBind(m.{{.M.SchemaName}}, echo.ExcludeFieldName(`created`))
		if err == nil {
			m.{{.M.IDField}} = {{.M.IDField|LowerCaseFirst}}
			err = m.Edit(nil, db.Cond{`{{.M.IDColumn}}`: {{.M.IDField|LowerCaseFirst}}})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/{{.Group}}/{{.H.TmplName "index"}}`))
			}
		}
	} {{if and .M.IDColumn (.M.HasAnyFields .M.SwitchableFields)}}else if ctx.IsAjax() {
		{{- range $k, $structField := .M.SwitchableFields -}}
		{{- $tableField := SnakeCase $structField -}}
		{{- if $.M.HasAnyField $structField}}
		{{$tableField}} := ctx.Query(`{{$tableField}}`)
		if len({{$tableField}}) > 0 {
			m.{{$structField}} = {{$tableField}}
			data := ctx.Data()
			err = m.SetField(nil, `{{$tableField}}`, {{$tableField}}, db.Cond{`{{$.M.IDColumn}}`: {{$.M.IDField|LowerCaseFirst}}})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
		{{end -}}
		{{- end -}}
	}{{end}} else if err == nil {
		echo.StructToForm(ctx, m.{{.M.SchemaName}}, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/{{.Group}}/{{.H.TmplName "index"}}`)
	ctx.Set(`title`, ctx.T(`编辑%v`, {{.M.PkgName}}.ObjectName{{.M.Name}}))
	ctx.Set(`listPageTitle`, ctx.T(`%v列表`, {{.M.PkgName}}.ObjectName{{.M.Name}}))
	return ctx.Render(`{{.Group}}/{{.H.TmplName "edit"}}`, err)
}

func {{.H.Name}}Delete(ctx echo.Context) error {
	{{.M.IDField|LowerCaseFirst}} := ctx.FormxValues(`{{.M.IDField|LowerCaseFirst}}`).{{.M.IDFieldType|Title}}(func(index int, value {{.M.IDFieldType}}) bool {
		return value > 0
	})
	m := {{.M.PkgName}}.New{{.M.Name}}(ctx)
	var err error
	for _, _v := range {{.M.IDField|LowerCaseFirst}} {
		if err = m.Delete(nil, db.Cond{`{{.M.IDColumn}}`: _v}); err != nil {
			break
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/{{.Group}}/{{.H.TmplName "index"}}`))
}
