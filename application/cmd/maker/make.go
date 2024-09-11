package maker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/coscms/webcore/library/config"
)

// Make 生成代码
func Make(c *CLIConfig) error {
	dbKey := c.DBKey
	if len(dbKey) == 0 {
		dbKey = factory.DefaultDBKey
	}
	fields := factory.DBIGet(dbKey).Fields
	group := c.Group
	tablePrefix := strings.Replace(group, `/`, `_`, -1)
	if len(tablePrefix) > 0 {
		tablePrefix += `_`
	}
	tables := c.Tables
	route := NewRoute()
	tplOpts := echo.H{`maker`: true}
	pkgName := path.Base(group)
	var tableList []string
	var err error
	if len(tables) == 0 { // 没有指定表信息的时候，自动根据表前缀去数据库查
		var allTables []string
		allTables, err = mysql.GetTables(0)
		if err != nil {
			return err
		}
		for _, table := range allTables {
			if strings.HasPrefix(table, tablePrefix) {
				tableList = append(tableList, strings.TrimPrefix(table, tablePrefix))
			}
		}
	} else {
		tableList = strings.Split(tables, `,`)
	}
	var switchableFieldList []string
	switchableFields := c.SwitchableFields
	if len(switchableFields) > 0 {
		switchableFieldList = strings.Split(switchableFields, ",")
		for index, field := range switchableFieldList {
			switchableFieldList[index] = com.PascalCase(field)
		}
	}
	for _, table := range tableList {
		table = strings.TrimSpace(table)
		if len(table) == 0 {
			continue
		}
		var objectName string
		objects := strings.SplitN(table, `:`, 2)
		switch len(objects) {
		case 2:
			objectName = strings.TrimSpace(objects[1])
			fallthrough
		case 1:
			table = strings.TrimSpace(objects[0])
			objectName, err = mysql.TableComment(0, config.FromFile().DB.Database, tablePrefix+table)
			if err != nil {
				return err
			}
		}
		fullTable := tablePrefix + table
		fieldList, ok := fields[fullTable]
		if !ok {
			fmt.Println(`Not Found Table:`, fullTable)
			continue
		}
		cfg := &Config{Group: group}
		cfg.H.Name = com.PascalCase(table)
		cfg.M.Name = cfg.H.Name
		cfg.M.Object = objectName
		cfg.M.PkgName = pkgName
		cfg.M.SchemaName = com.PascalCase(fullTable)
		cfg.M.Database = config.FromFile().DB.Database
		cfg.M.DBKey = dbKey
		cfg.M.NameField = ``
		cfg.M.NameColumn = ``
		cfg.M.IDField = ``
		cfg.M.IDFieldType = ``
		cfg.M.IDColumn = ``
		if len(switchableFieldList) > 0 {
			cfg.M.SwitchableFields = append(cfg.M.SwitchableFields, switchableFieldList...)
		}
		cfg.T.Options = tplOpts
		if info, ok := fieldList[`id`]; ok {
			cfg.M.IDField = `Id`
			cfg.M.IDFieldType = info.GoType
			cfg.M.IDColumn = `id`
		} else {
			for _, info := range fieldList {
				if info.PrimaryKey {
					cfg.M.IDField = info.GoName
					cfg.M.IDFieldType = info.GoType
					cfg.M.IDColumn = info.Name
					break
				}
			}
			if len(cfg.M.IDField) == 0 && !c.MustHasPrimaryKey {
				cfg.M.IDField = `Id`
				cfg.M.IDFieldType = `uint64`
				cfg.M.IDColumn = `id`
			}
		}
		for _, nameCol := range []string{`name`, `title`} {
			if info, ok := fieldList[nameCol]; ok {
				cfg.M.NameField = info.GoName
				cfg.M.NameColumn = info.Name
			}
		}
		if err := MakeModel(cfg); err != nil {
			return err
		}
		if len(cfg.M.IDField) > 0 {
			if err := MakeHandler(cfg); err != nil {
				return err
			}
			if err := MakeTemplate(cfg); err != nil {
				return err
			}
			cfg.MakeInit(route)
		}
	}
	initCode := route.String()
	initNavigation := route.NavString()
	data := echo.H{
		`MakeInit`:           initCode,
		`MakeInitNavigation`: initNavigation,
		`PkgName`:            pkgName,
		`Group`:              group,
	}
	return MakeHandlerInit(group, data)
}

func compile(tmpl string, data interface{}) ([]byte, error) {
	t := template.New(filepath.Base(tmpl))
	b, err := os.ReadFile(tmpl)
	if err != nil {
		return nil, err
	}
	t.Funcs(template.FuncMap(tplfunc.TplFuncMap))
	_, err = t.Parse(com.Bytes2str(b))
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(``)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func Format(file string) error {
	cmd := exec.Command(`gofmt`, `-l`, `-s`, `-w`, file)
	return cmd.Run()
}
