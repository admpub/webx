//go:build !bindata
// +build !bindata

package xtemplate

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/webx-top/echo"
)

func (t *Template) Fix(ctx echo.Context, _ http.FileSystem, theme string, tmpl string) (string, bool) {
	return t.PathFixers.Fix(ctx, t, theme, tmpl)
}

// Fix 模版路径修复
func (p *PathFixers) Fix(ctx echo.Context, t *Template, theme string, tmpl string) (string, bool) {
	cacheKey := theme + `>` + tmpl
	if mp, ok := t.cachedPathData.get(cacheKey); ok {
		return mp.String, mp.Valid
	}
	subdir, newTmpl, group := p.parsePath(theme, tmpl)
	if len(group) > 0 {
		if t, ok := groups[group]; ok {
			r := t.Handle(ctx, subdir, tmpl)
			t.cachedPathData.set(cacheKey, sql.NullString{String: r, Valid: true})
			return r, true
		}
	}
	pathFixers, ok := (*p)[subdir]
	if ok {
		if _tmpl, ok := findPath(pathFixers, subdir, newTmpl); ok {
			t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
			return _tmpl, ok
		}
	}
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
			pathFixers, ok := (*p)[fb]
			if ok {
				subdir = fb
				if _tmpl, ok := findPath(pathFixers, subdir, newTmpl); ok {
					t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
					return _tmpl, ok
				}
			}
			_tmpl := filepath.Join(t.TmplDir, fb, newTmpl)
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

func findPath(pathFixers []PathFixer, subdir string, newTmpl string) (string, bool) {
	for _, pathFixer := range pathFixers {
		_tmpl := pathFixer(subdir, newTmpl)
		fi, err := os.Stat(_tmpl)
		if err == nil && !fi.IsDir() {
			return _tmpl, true
		}
	}
	return ``, false
}
