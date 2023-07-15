package route

import (
	"strings"

	"github.com/webx-top/echo"
)

func NewGroup(groupNamers ...func(string) string) *Group {
	return &Group{
		Handlers:       map[string][]func(echo.RouteRegister){},
		Middlewares:    map[string][]interface{}{},
		PreMiddlewares: map[string][]interface{}{},
		Namers:         groupNamers,
		Metas:          map[string]echo.H{},
	}
}

type Group struct {
	Handlers       map[string][]func(echo.RouteRegister)
	Namers         []func(string) string
	Middlewares    map[string][]interface{}
	PreMiddlewares map[string][]interface{}
	Metas          map[string]echo.H
}

func (g *Group) AddNamer(namers ...func(string) string) {
	g.Namers = append(g.Namers, namers...)
}

func (g *Group) SetNamer(namers ...func(string) string) {
	g.Namers = namers
}

func (g *Group) Pre(groupName string, middlewares ...interface{}) {
	if _, ok := g.PreMiddlewares[groupName]; !ok {
		g.PreMiddlewares[groupName] = []interface{}{}
	}
	g.PreMiddlewares[groupName] = append(g.PreMiddlewares[groupName], middlewares...)
}

func (g *Group) Use(groupName string, middlewares ...interface{}) {
	if groupName != `*` && strings.HasSuffix(groupName, `*`) {
		groupName = strings.TrimRight(groupName, `*`)
		if _, ok := g.PreMiddlewares[groupName]; !ok {
			g.PreMiddlewares[groupName] = []interface{}{}
		}
		g.PreMiddlewares[groupName] = append(g.PreMiddlewares[groupName], middlewares...)
		return
	}
	if _, ok := g.Middlewares[groupName]; !ok {
		g.Middlewares[groupName] = []interface{}{}
	}
	g.Middlewares[groupName] = append(g.Middlewares[groupName], middlewares...)
}

func (g *Group) Register(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	_, ok := g.Handlers[groupName]
	if !ok {
		g.Handlers[groupName] = []func(echo.RouteRegister){}
	}
	if len(middlewares) > 0 {
		g.Use(groupName, middlewares...)
	}
	g.Handlers[groupName] = append(g.Handlers[groupName], fn)
}

func (g *Group) SetMeta(groupName string, meta echo.H) {
	if _, ok := g.Metas[groupName]; ok {
		g.Metas[groupName].DeepMerge(meta)
	} else {
		g.Metas[groupName] = meta
	}
}

func (g *Group) SetMetaKV(groupName string, key string, val interface{}) {
	if _, ok := g.Metas[groupName]; !ok {
		g.Metas[groupName] = echo.H{}
	}
	g.Metas[groupName].Set(key, val)
}

func (g *Group) Apply(e echo.RouteRegister, rootGroup string) {

	var groupDefaultMiddlewares []interface{}
	middlewares, ok := g.Middlewares[`*`]
	if ok {
		groupDefaultMiddlewares = append(groupDefaultMiddlewares, middlewares...)
	}
	for group, handlers := range g.Handlers {
		originalGroupName := group
		for _, namer := range g.Namers {
			group = namer(group)
		}
		if len(group) > 0 && !strings.HasPrefix(group, `/`) {
			panic(`Routing group name must start with "/". Your name is passed: ` + group)
		}
		grp := e.Group(group)
		if g.Metas != nil {
			if meta, ok := g.Metas[originalGroupName]; ok {
				grp.SetMeta(meta)
			}
		}
		if group != rootGroup { // 组名为空时，为顶层组
			grp.Use(groupDefaultMiddlewares...)
		}
		for prefix, middlewares := range g.PreMiddlewares {
			if strings.HasPrefix(originalGroupName, prefix) {
				grp.Use(middlewares...)
			}
		}
		middlewares, ok := g.Middlewares[originalGroupName]
		if ok {
			grp.Use(middlewares...)
		}
		for _, register := range handlers {
			register(grp)
		}
	}
}
