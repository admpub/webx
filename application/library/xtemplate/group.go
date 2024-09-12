package xtemplate

import (
	"database/sql"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render/driver"
)

func New(kind string, pa *ntemplate.PathAliases, registerToGroup ...bool) *Template {
	if pa == nil {
		pa = ntemplate.NewPathAliases()
	}
	t := &Template{
		Kind:            kind,
		PathFixers:      &PathFixers{},
		PathAliases:     pa,
		themeInfo:       &ThemeInfo{Name: `default`},
		themeMutex:      sync.RWMutex{},
		themeOnce:       sync.Once{},
		cachedPathData:  newCachedPathData(),
		themeInfoStorer: NewFileStore(kind),
	}
	if len(registerToGroup) > 0 && registerToGroup[0] {
		Register(kind, t)
	}
	return t
}

func newCachedPathData() *cachedPathData {
	return &cachedPathData{
		mapping: map[string]sql.NullString{},
		mutex:   sync.RWMutex{},
	}
}

type cachedPathData struct {
	mapping map[string]sql.NullString
	mutex   sync.RWMutex
}

type Template struct {
	handler     PathHandle
	TmplDir     string
	Kind        string
	PathAliases *ntemplate.PathAliases
	customFS    http.FileSystem
	enableTheme bool
	*PathFixers
	themeInfo       *ThemeInfo
	themeMutex      sync.RWMutex
	themeOnce       sync.Once
	cachedPathData  *cachedPathData
	themeInfoStorer Storer
}

func (t *Template) SetTmplDir(tmplDir string) *Template {
	t.TmplDir = tmplDir
	return t
}

func (t *Template) SetStorer(storer Storer) *Template {
	t.themeInfoStorer = storer
	return t
}

func (t *Template) SetCustomFS(fs http.FileSystem) *Template {
	t.customFS = fs
	return t
}

func (t *Template) SetEnableTheme(enable bool) *Template {
	t.enableTheme = enable
	return t
}

func (t *Template) Storer() Storer {
	return t.themeInfoStorer
}

func (t *Template) SetHandler(h PathHandle) *Template {
	t.handler = h
	return t
}

func (t *Template) Handle(ctx echo.Context, theme string, tmpl string) string {
	return t.handler(ctx, theme, tmpl)
}

func (t *Template) SetPathFixers(h *PathFixers) *Template {
	t.PathFixers = h
	return t
}

func (t *Template) ApplyAliases() *Template {
	t.PathAliases.Range(func(prefix, templateDir string) error {
		log.Debug(`[`+t.Kind+`] `, `Template subfolder "`+prefix+`" is relocated to: `, templateDir)
		t.AddDir(prefix, templateDir)
		return nil
	})
	return t
}

func (t *Template) AddAlias(alias, tmplDir string) *Template {
	t.PathAliases.Add(alias, tmplDir)
	return t
}

func (t *Template) AddAliasFromAllSubdir(tmplDir string) *Template {
	t.PathAliases.AddAllSubdir(tmplDir)
	return t
}

func (t *Template) TmplDirs() []string {
	return t.PathAliases.TmplDirs()
}

func (t *Template) Aliases() []string {
	return t.PathAliases.Aliases()
}

func (t *Template) Register(renderer driver.Driver, watchOtherDirs ...string) {
	hasCustomFS := t.customFS != nil
	// 设置后台模板路径修正器的修正处理函数
	t.SetTmplDir(renderer.TmplDir()).SetHandler(func(c echo.Context, theme string, tmpl string) string {
		var found bool
		tmpl, found = t.Fix(c, t.customFS, theme, tmpl)
		if found {
			return tmpl
		}
		if hasCustomFS {
			return path.Join(t.Kind, tmpl)
		}
		if t.enableTheme && len(theme) > 0 {
			tmpl = theme + `/` + tmpl
		}
		return filepath.Join(renderer.TmplDir(), tmpl)
	})

	// 将后台模板路径修正器与模板渲染引擎关联
	renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
		var theme string
		if t.enableTheme {
			theme = c.Internal().String(`theme`, `default`)
		}
		return t.Handle(c, theme, tmpl)
	})

	// 关注后台模板路径内的文件改动
	if !hasCustomFS {
		for _, watchOtherDir := range watchOtherDirs {
			if len(watchOtherDir) > 0 {
				renderer.Manager().AddWatchDir(watchOtherDir)
			}
		}
		for _, templateDir := range t.PathAliases.TmplDirs() {
			log.Debug(`[`+t.Kind+`] `, `Watch folder: `, templateDir)
			renderer.Manager().AddWatchDir(templateDir)
		}
	}
}

func (t *Template) ThemeInfo(c echo.Context) *ThemeInfo {
	t.themeMutex.RLock()
	t.themeOnce.Do(func() {
		t.loadThemeInfo(c)
	})
	v := t.themeInfo
	t.themeMutex.RUnlock()
	return v
}

func (t *Template) loadThemeInfo(c echo.Context) {
	var err error
	t.themeInfo, err = t.themeInfoStorer.Get(c, ``)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Error(err)
		return
	}
}

func (t *Template) SetThemeInfo(c echo.Context, v *ThemeInfo) {
	t.themeMutex.Lock()
	t.themeInfo = v
	t.themeInfoStorer.Put(c, ``, v)
	t.cachedPathData.clear()
	t.themeMutex.Unlock()
}

func (t *Template) ClearCache() {
	t.cachedPathData.clear()
}

func (t *cachedPathData) get(path string) (sql.NullString, bool) {
	t.mutex.RLock()
	mp, ok := t.mapping[path]
	t.mutex.RUnlock()
	return mp, ok
}

func (t *cachedPathData) clear() {
	t.mutex.Lock()
	t.mapping = map[string]sql.NullString{}
	t.mutex.Unlock()
}

func (t *cachedPathData) set(path string, mp sql.NullString) {
	t.mutex.Lock()
	t.mapping[path] = mp
	t.mutex.Unlock()
}

var groups = map[string]*Template{}

func Register(group string, t *Template) {
	groups[group] = t
}

func Unregister(group string) {
	delete(groups, group)
}

func Get(group string) *Template {
	if v, ok := groups[group]; ok {
		return v
	}
	return nil
}

func Backend() *Template {
	return Get(KindBackend)
}

func Frontend() *Template {
	return Get(KindFrontend)
}
