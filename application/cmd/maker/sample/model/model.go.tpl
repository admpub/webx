package {{.M.PkgName}}

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

const ObjectName{{.M.Name}} = `{{.M.Object}}`

func New{{.M.Name}}(ctx echo.Context) *{{.M.Name}} {
	m := &{{.M.Name}}{
		{{.M.SchemaName}}: dbschema.New{{.M.SchemaName}}(ctx),
	}
	return m
}

type {{.M.Name}} struct {
	*dbschema.{{.M.SchemaName}}
}
{{if .M.NameColumn}}
func (f *{{.M.Name}}) Exists(name string) (bool, error) {
	return f.{{.M.SchemaName}}.Exists(nil, db.Cond{`{{.M.NameColumn}}`: name})
}

func (f *{{.M.Name}}) ExistsOther(name string, {{.M.IDField|LowerCaseFirst}} {{.M.IDFieldType}}) (bool, error) {
	return f.{{.M.SchemaName}}.Exists(nil, db.And(
		db.Cond{`{{.M.NameColumn}}`: name},
		db.Cond{`{{.M.IDColumn}}`: db.NotEq({{.M.IDField|LowerCaseFirst}})},
	))
}
{{end}}
func (f *{{.M.Name}}) check() error {
	{{if .M.NameColumn}}
	var (
		exists bool
		err error
	)
	if f.{{.M.IDField}} < 1 {
		exists, err = f.Exists(f.{{.M.NameField}})
	} else {
		exists, err = f.ExistsOther(f.{{.M.NameField}}, f.{{.M.IDField}})
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().E(`名称“%s”已存在`, f.{{.M.NameField}})
	}
	{{end}}
	return nil
}

func (f *{{.M.Name}}) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.{{.M.SchemaName}}.Insert()
}

func (f *{{.M.Name}}) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.{{.M.SchemaName}}.Update(mw, args...)
}

func (f *{{.M.Name}}) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*dbschema.{{.M.SchemaName}}, error) {
	list := []*dbschema.{{.M.SchemaName}}{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	return list, err
}
{{if .M.IDField}}
func (f *{{.M.Name}}) NextRow(currentID {{.M.IDFieldType}}, extraCond *db.Compounds) (*dbschema.{{.M.SchemaName}}, error) {
	row := &dbschema.{{.M.SchemaName}}{}
	row.CPAFrom(f.{{.M.SchemaName}})
	cond := db.NewCompounds()
	//cond.AddKV(`disabled`, `N`)
	cond.AddKV(`{{.M.IDField}}`, db.Lt(currentID))
	if extraCond != nil {
		cond.From(extraCond)
	}
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`*`).OrderBy(`-{{.M.IDField}}`)
	}, cond.And())
	return row, err
}

func (f *{{.M.Name}}) PrevRow(currentID {{.M.IDFieldType}}, extraCond *db.Compounds) (*dbschema.{{.M.SchemaName}}, error) {
	row := &dbschema.{{.M.SchemaName}}{}
	row.CPAFrom(f.{{.M.SchemaName}})
	cond := db.NewCompounds()
	//cond.AddKV(`disabled`, `N`)
	cond.AddKV(`{{.M.IDField}}`, db.Gt(currentID))
	if extraCond != nil {
		cond.From(extraCond)
	}
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`*`).OrderBy(`{{.M.IDField}}`)
	}, cond.And())
	return row, err
}
{{end}}