package maker

import (
	"fmt"
	"path"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/registry/navigate"
)

// HandlerConfig Handler配置
type HandlerConfig struct {
	Name string //Handler名称
}

// FileName Handler文件名
func (h *HandlerConfig) FileName() string {
	return com.SnakeCase(h.Name)
}

// TmplName 模板名
func (h *HandlerConfig) TmplName(typ string) string {
	if len(h.Name) == 0 {
		return typ
	}
	return h.FileName() + `_` + typ
}

// ModelConfig Model配置
type ModelConfig struct {
	Name             string               //Model名称(首字母大写)
	Object           string               //管理的目标对象名称(如:新闻、商品等)
	PkgName          string               //Model包名称
	SchemaName       string               //Schema名称(dbschema内的结构体名称)
	NameField        string               //名称字段名(结构体)
	NameColumn       string               //名称列名(数据库)
	IDField          string               //ID字段名(结构体)
	IDFieldType      string               //ID字段数据类型(结构体)
	IDColumn         string               //ID列名(数据库)
	fields           []*factory.FieldInfo //字段信息
	Database         string               //数据库名称
	DBKey            string               //数据库标识
	SwitchableFields []string             //可切换状态的字段（即类型为枚举值enum('Y','N')的字段）
}

var emptyFields = []*factory.FieldInfo{}

// FileName Model文件名
func (m *ModelConfig) FileName() string {
	return com.SnakeCase(m.Name)
}

// Fields 字段信息
func (m *ModelConfig) Fields() []*factory.FieldInfo {
	if m.fields != nil {
		return m.fields
	}
	tableName := com.SnakeCase(m.SchemaName)
	fields := make([]*factory.FieldInfo, 0)
	dbi := factory.DBIGet(m.DBKey)
	f := dbi.Fields
	columns := dbi.TableColumns(tableName)
	for _, field := range columns {
		if info, exists := f.Find(tableName, field); exists {
			fields = append(fields, info)
		}
	}
	m.fields = fields
	return m.fields
}

// HasAnyField 包含任意一个字段名
func (m *ModelConfig) HasAnyField(fields ...string) bool {
	return m.HasAnyFields(fields)
}

// HasAnyFields 包含任意一个字段名
func (m *ModelConfig) HasAnyFields(fields []string) bool {
	infos := m.Fields()
	for _, field := range fields {
		for _, info := range infos {
			if info.GoName == field {
				return true
			}
		}
	}
	return false
}

// HasAnyColumn 包含任意一个列名
func (m *ModelConfig) HasAnyColumn(columns ...string) bool {
	return m.HasAnyColumns(columns)
}

// HasAnyColumns 包含任意一个列名
func (m *ModelConfig) HasAnyColumns(columns []string) bool {
	infos := m.Fields()
	for _, column := range columns {
		for _, info := range infos {
			if info.Name == column {
				return true
			}
		}
	}
	return false
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	Name    string //模板名称
	Options echo.H //模板选项
}

func NewConfig() *Config {
	return &Config{
		M: ModelConfig{
			SwitchableFields: []string{"Disabled", "Display"},
		},
	}
}

// Config 配置
type Config struct {
	Group string //组名称(请用全小写的英文字符) official/exmple
	H     HandlerConfig
	M     ModelConfig
	T     TemplateConfig
}

// PkgName Handler包名
func (c *Config) PkgName() string {
	if len(c.Group) == 0 {
		return `handler`
	}
	return path.Base(c.Group)
}

type Route struct {
	NavList  *[]navigate.Item
	Routes   []string
	GRoutes  map[string][]string
	G2Routes map[string]map[string][]string
}

func NewRoute() *Route {
	return &Route{
		NavList:  &[]navigate.Item{},
		Routes:   []string{},
		GRoutes:  map[string][]string{},
		G2Routes: map[string]map[string][]string{},
	}
}

func (r *Route) String() string {
	var s string
	var sr string
	for _, route := range r.Routes {
		sr += route
	}
	if len(sr) > 0 {
		s += "\n\t" + `route.Register(func(g echo.RouteRegister) {` + sr + "\t" + `})` + "\n"
	}
	sr = ``
	for group, routes := range r.GRoutes {
		var ss string
		for _, route := range routes {
			ss += route
		}
		groups, ok := r.G2Routes[group]
		if ok {
			for subGrp, subRts := range groups {
				ss += "\n\t\tg = g.Group(`/" + subGrp + "`)"
				for _, route := range subRts {
					ss += route
				}
			}
			delete(r.G2Routes, group)
		}
		sr += "\n\t" + `route.RegisterToGroup("/` + group + `", func(g echo.RouteRegister) {` + ss + "\t" + `})` + "\n"
	}
	for group, grpRoutes := range r.G2Routes {
		var ss string
		for subGrp, subRts := range grpRoutes {
			ss += "\n\t\tg = g.Group(`/" + subGrp + "`)"
			for _, route := range subRts {
				ss += route
			}
		}
		sr += "\n\t" + `route.RegisterToGroup("/` + group + `", func(g echo.RouteRegister) {` + ss + "\t" + `})` + "\n"
	}
	s += sr
	return s
}

func (r *Route) NavString() string {
	s := fmt.Sprintf("%#v", r.NavList)
	s = strings.ReplaceAll(s, `Children:(*navigate.List)(nil)`, ``)
	s = strings.ReplaceAll(s, `&[]navigate.Item{`, `&navigate.List{`)
	s = strings.ReplaceAll(s, `navigate.Item{`, `{`)
	s = strings.ReplaceAll(s, `, `, ",\n")
	s = strings.ReplaceAll(s, `{`, "{\n")
	s = strings.ReplaceAll(s, `}}`, "},\n}")
	return s
}

// MakeInit 生成init代码
func (c *Config) MakeInit(r *Route) *Route {
	*r.NavList = append(*r.NavList, navigate.Item{
		Display:   true,
		Name:      c.M.Object + `管理`,     //名称
		Action:    c.H.TmplName(`index`), //操作(一般为网址)
		Icon:      ``,                    //图标
		Target:    ``,                    //打开方式
		Unlimited: false,                 //是否不限制权限
	})
	*r.NavList = append(*r.NavList, navigate.Item{
		Display:   false,
		Name:      `添加` + c.M.Object,   //名称
		Action:    c.H.TmplName(`add`), //操作(一般为网址)
		Icon:      `plus`,              //图标
		Target:    ``,                  //打开方式
		Unlimited: false,               //是否不限制权限
	})
	*r.NavList = append(*r.NavList, navigate.Item{
		Display:   false,
		Name:      `修改` + c.M.Object,    //名称
		Action:    c.H.TmplName(`edit`), //操作(一般为网址)
		Icon:      `pencil`,             //图标
		Target:    ``,                   //打开方式
		Unlimited: false,                //是否不限制权限
	})
	*r.NavList = append(*r.NavList, navigate.Item{
		Display:   false,
		Name:      `删除` + c.M.Object,      //名称
		Action:    c.H.TmplName(`delete`), //操作(一般为网址)
		Icon:      `remove`,               //图标
		Target:    ``,                     //打开方式
		Unlimited: false,                  //是否不限制权限
	})
	routeReg := `
		g.Route("GET,POST", "/` + c.H.TmplName(`index`) + `", ` + c.H.Name + `Index)
		g.Route("GET,POST", "/` + c.H.TmplName(`add`) + `", ` + c.H.Name + `Add)
		g.Route("GET,POST", "/` + c.H.TmplName(`edit`) + `", ` + c.H.Name + `Edit)
		g.Route("GET,POST", "/` + c.H.TmplName(`delete`) + `", ` + c.H.Name + `Delete)
`
	if len(c.Group) == 0 {
		r.Routes = append(r.Routes, routeReg)
		return r
	}
	var group, route string
	inf := strings.SplitN(c.Group, `/`, 2)
	switch len(inf) {
	case 2:
		route = inf[1]
		group = inf[0]
		if _, ok := r.G2Routes[group]; !ok {
			r.G2Routes[group] = map[string][]string{}
		}
		if _, ok := r.G2Routes[group][route]; !ok {
			r.G2Routes[group][route] = []string{}
		}
		r.G2Routes[group][route] = append(r.G2Routes[group][route], routeReg)
	case 1:
		group = inf[0]
		if _, ok := r.GRoutes[group]; !ok {
			r.GRoutes[group] = []string{}
		}
		r.GRoutes[group] = append(r.GRoutes[group], routeReg)
	}
	return r
}
