//go:build bindata
// +build bindata

package xtemplate

import (
	"database/sql"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/webx-top/echo"
)

func (t *Template) Fix(ctx echo.Context, fs http.FileSystem, theme string, tmpl string) (string, bool) {
	return t.PathFixers.Fix(ctx, fs, t, theme, tmpl)
}

// Fix 模版路径修复
func (p *PathFixers) Fix(ctx echo.Context, fs http.FileSystem, t *Template, theme string, tmpl string) (string, bool) {
	cacheKey := theme + `>` + tmpl
	if mp, ok := t.cachedPathData.get(cacheKey); ok {
		return mp.String, mp.Valid
	}
	subdir, tplfile, group := p.parsePath(theme, tmpl)
	if len(group) > 0 {
		if t, ok := groups[group]; ok {
			r := t.Handle(ctx, subdir, tmpl)
			t.cachedPathData.set(cacheKey, sql.NullString{String: r, Valid: true})
			return r, true
		}
	}

	dirName := filepath.Base(t.TmplDir)

	if _tmpl, exists := fsFileExists(fs, t.TmplDir, dirName, subdir, tplfile); exists {
		t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
		return _tmpl, exists
	}

	if themeInfo := GetThemeInfoFromContext(ctx); themeInfo != nil && len(themeInfo.Fallback) > 0 {
		rawTmpl := tplfile
		for _, fb := range themeInfo.Fallback {
			if len(fb) == 0 {
				continue
			}

			if _tmpl, exists := fsFileExists(fs, t.TmplDir, dirName, fb, rawTmpl); exists {
				t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
				return _tmpl, exists
			}
		}
	}
	t.cachedPathData.set(cacheKey, sql.NullString{String: tmpl})
	return tmpl, false
}

func fsFileExists(fs http.FileSystem, tmplDir, dirName, subDir, tmpl string) (string, bool) {
	_tmpl := filepath.Join(tmplDir, subDir, tmpl)
	fi, err := os.Stat(_tmpl)
	if err == nil && !fi.IsDir() {
		_tmpl = path.Join(dirName, subDir, tmpl)
		return _tmpl, true
	}
	_tmpl = path.Join(dirName, subDir, tmpl)
	file, err := fs.Open(_tmpl)
	if err == nil {
		var fi os.FileInfo
		fi, err = file.Stat()
		file.Close()
		if err == nil && !fi.IsDir() {
			return _tmpl, true
		}
	}
	return ``, false
}
