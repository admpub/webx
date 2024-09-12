// 模板文件系统、模板文件路径修正器、模板主题配置
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

func (p *PathFixers) parsePath(theme string, tmpl string) (subdir string, newTmpl string, group string) {
	if len(tmpl) > 3 && strings.HasPrefix(tmpl, startChar) { // #theme#index
		subdir = tmpl[1:]
		splited := strings.SplitN(subdir, endChar, 2)
		subdir = splited[0]
		if len(splited) > 1 {
			tmpl = strings.TrimLeft(splited[1], whitespace)
			subdirAndGroup := strings.SplitN(subdir, groupAtChar, 2) // #theme@group#
			if len(subdirAndGroup) == 2 {
				subdir = strings.TrimRight(subdirAndGroup[0], whitespace)
				group = strings.TrimLeft(subdirAndGroup[1], whitespace)
				return
			}
		}
		if len(theme) == 0 { // 未启用主题
			subdir, newTmpl = p.splitPath(tmpl)
		} else {
			newTmpl = tmpl
		}
	} else if len(theme) > 0 {
		subdir = theme
		newTmpl = tmpl
	} else {
		subdir, newTmpl = p.splitPath(tmpl)
	}
	return
}

const (
	startChar   = `#`
	endChar     = `#`
	groupAtChar = `@`
	whitespace  = ` `
)
