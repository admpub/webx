package xtemplate

import (
	"database/sql"
	"os"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

func New(kind string) *Template {
	return &Template{
		Kind:            kind,
		PathFixers:      &PathFixers{},
		themeInfo:       &ThemeInfo{Name: `default`},
		themeMutex:      sync.RWMutex{},
		themeOnce:       sync.Once{},
		cachedPathData:  newCachedPathData(),
		themeInfoStorer: NewFileStore(kind),
	}
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
	handler PathHandle
	TmplDir string
	Kind    string
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
