//go:build bindata
// +build bindata

package xtemplate

import (
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"strings"

	assetfs "github.com/admpub/go-bindata-assetfs"
	"github.com/webx-top/echo"
)

func (t *Template) Fix(ctx echo.Context, fs *assetfs.AssetFS, theme string, tmpl string) (string, bool) {
	return t.PathFixers.Fix(ctx, fs, t, theme, tmpl)
}

// Fix 模版路径修复
func (p *PathFixers) Fix(ctx echo.Context, fs *assetfs.AssetFS, t *Template, theme string, tmpl string) (string, bool) {
	cacheKey := theme + `>` + tmpl
	if mp, ok := t.cachedPathData.get(cacheKey); ok {
		return mp.String, mp.Valid
	}
	var subdir string
	tplfile := tmpl
	if len(tmpl) > 3 && strings.HasPrefix(tmpl, startChar) { // #theme#index
		subdir = tmpl[1:]                             // theme#index
		splited := strings.SplitN(subdir, endChar, 2) //[theme,index]
		subdir = splited[0]
		if len(splited) > 1 {
			tplfile = strings.TrimLeft(splited[1], whitespace)
			subdirAndGroup := strings.SplitN(subdir, groupAtChar, 2) // #theme@group#
			if len(subdirAndGroup) == 2 {
				subdir = strings.TrimRight(subdirAndGroup[0], whitespace)
				group := strings.TrimLeft(subdirAndGroup[1], whitespace)
				if t, ok := groups[group]; ok {
					r := t.Handle(ctx, subdir, tplfile)
					t.cachedPathData.set(cacheKey, sql.NullString{String: r, Valid: true})
					return r, true
				}
			}
		}
		//echo.Dump(echo.H{`subdir`: subdir, `tplfile`: tplfile, `theme`: theme})
		if len(theme) > 0 { // 启用主题
			tmpl = subdir + `/` + tplfile
		} else {
			tmpl = tplfile
		}
	} else if len(theme) > 0 {
		tmpl = theme + `/` + tplfile
	}
	_tmpl := filepath.Join(t.TmplDir, tmpl)
	fi, err := os.Stat(_tmpl)
	if err == nil && !fi.IsDir() {
		_tmpl = path.Join(filepath.Base(t.TmplDir), tmpl)
		t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
		return _tmpl, true
	}
	file, err := fs.Open(tmpl)
	if err == nil {
		var fi os.FileInfo
		fi, err = file.Stat()
		file.Close()
		if err == nil && !fi.IsDir() {
			_tmpl = path.Join(filepath.Base(t.TmplDir), tmpl)
			t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
			return _tmpl, true
		}
	}
	if themeInfo := GetThemeInfoFromContext(ctx); themeInfo != nil && len(themeInfo.Fallback) > 0 {
		rawTmpl := tplfile
		dirName := filepath.Base(t.TmplDir)
		for _, fb := range themeInfo.Fallback {
			if len(fb) == 0 {
				continue
			}
			_tmpl := filepath.Join(t.TmplDir, fb, rawTmpl)
			fi, err := os.Stat(_tmpl)
			if err == nil && !fi.IsDir() {
				_tmpl = path.Join(dirName, fb, rawTmpl)
				t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
				return _tmpl, true
			}
			_tmpl = path.Join(dirName, fb, rawTmpl)
			file, err = fs.Open(_tmpl)
			if err == nil {
				var fi os.FileInfo
				fi, err = file.Stat()
				file.Close()
				if err == nil && !fi.IsDir() {
					t.cachedPathData.set(cacheKey, sql.NullString{String: _tmpl, Valid: true})
					return _tmpl, true
				}
			}
		}
	}
	t.cachedPathData.set(cacheKey, sql.NullString{String: tmpl})
	return tmpl, false
}
