package xtemplate

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/webx-top/echo"
)

type PathFixer func(subdir, tmpl string) string
type PathHandle func(c echo.Context, theme string, tmpl string) string
type PathFixers map[string][]PathFixer //模版路径 {subdir:func}

func (p *PathFixers) Add(dirName string, fixer PathFixer) *PathFixers {
	if _, ok := (*p)[dirName]; !ok {
		(*p)[dirName] = []PathFixer{}
	}
	(*p)[dirName] = append((*p)[dirName], fixer)
	return p
}

func (p *PathFixers) AddDir(dirName string, parentPath string) *PathFixers {
	p.Add(dirName, p.MakeFixer(parentPath))
	return p
}

func (p *PathFixers) MakeFixer(parentPath string) PathFixer {
	if !filepath.IsAbs(parentPath) {
		var err error
		parentPath, err = filepath.Abs(parentPath)
		if err != nil {
			panic(err)
		}
	}
	return func(subdir, tmpl string) string {
		return filepath.Join(parentPath, subdir, tmpl)
	}
}

func (p *PathFixers) Delete(names ...string) *PathFixers {
	for _, name := range names {
		delete(*p, name)
	}
	return p
}

func (p *PathFixers) Keys() []string {
	names := make([]string, len(*p))
	var i int
	for name := range *p {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

func (p *PathFixers) splitPath(tmpl string) (subdir string, newTmpl string) {
	newTmpl = tmpl
	pos := strings.Index(tmpl, `/`)
	if pos < 0 {
		return
	}
	subdir = tmpl[:pos]
	if len(tmpl) > pos+1 {
		newTmpl = tmpl[pos+1:]
	} else {
		newTmpl = ``
	}
	return
}

const (
	startChar   = `#`
	endChar     = `#`
	groupAtChar = `@`
	whitespace  = ` `
)
