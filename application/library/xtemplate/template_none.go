//go:build !bindata
// +build !bindata

package xtemplate

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
)

func (t *Template) Fix(ctx echo.Context, theme string, tmpl string) (string, bool) {
	return t.PathFixers.Fix(ctx, t, theme, tmpl)
}

// Fix 模版路径修复
func (p *PathFixers) Fix(ctx echo.Context, t *Template, theme string, tmpl string) (string, bool) {
	cacheKey := theme + `>` + tmpl
	if mp, ok := t.cachedPathData.get(cacheKey); ok {
		return mp.String, mp.Valid
	}
	var subdir, newTmpl string
	if len(tmpl) > 3 && strings.HasPrefix(tmpl, startChar) { // #theme#index
		subdir = tmpl[1:]
		splited := strings.SplitN(subdir, endChar, 2)
		subdir = splited[0]
		if len(splited) > 1 {
			tmpl = strings.TrimLeft(splited[1], whitespace)
			subdirAndGroup := strings.SplitN(subdir, groupAtChar, 2) // #theme@group#
			if len(subdirAndGroup) == 2 {
				subdir = strings.TrimRight(subdirAndGroup[0], whitespace)
				group := strings.TrimLeft(subdirAndGroup[1], whitespace)
				if t, ok := groups[group]; ok {
					r := t.Handle(ctx, subdir, tmpl)
					t.cachedPathData.set(cacheKey, sql.NullString{String: r, Valid: true})
					return r, true
				}
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
	pathFixers, ok := (*p)[subdir]
	if !ok {
		_tmpl := filepath.Join(t.TmplDir, subdir, newTmpl)
		fi, err := os.Stat(_tmpl)
		if err == nil && !fi.IsDir() {
			t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
			return _tmpl, true
		}
		if themeInfo := GetThemeInfoFromContext(ctx); themeInfo != nil && len(themeInfo.Fallback) > 0 {
			for _, fb := range themeInfo.Fallback {
				if len(fb) == 0 {
					continue
				}
				pathFixers, ok = (*p)[fb]
				if ok {
					subdir = fb
					break
				}
				_tmpl := filepath.Join(t.TmplDir, fb, newTmpl)
				fi, err := os.Stat(_tmpl)
				if err == nil && !fi.IsDir() {
					t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
					return _tmpl, true
				}
			}
		}
	}
	if ok {
		if len(pathFixers) == 1 {
			r := pathFixers[0](subdir, newTmpl)
			t.cachedPathData.set(cacheKey, sql.NullString{String: r, Valid: true})
			return r, true
		}
		for _, pathFixer := range pathFixers {
			_tmpl := pathFixer(subdir, newTmpl)
			fi, err := os.Stat(_tmpl)
			if err == nil && !fi.IsDir() {
				t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
				return _tmpl, true
			}
		}
	}
	t.cachedPathData.set(cacheKey, sql.NullString{String: tmpl})
	return tmpl, false
}
