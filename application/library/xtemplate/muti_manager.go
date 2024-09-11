package xtemplate

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/bindata"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/manager"
)

func NewMultiManager(templateDir string, managers ...driver.Manager) driver.Manager {
	templateDir, _ = filepath.Abs(templateDir)
	var hasBindata bool
	for _, mgr := range managers {
		if _, ok := mgr.(*bindata.TmplManager); ok {
			hasBindata = true
			break
		}
	}
	return &MultiManager{
		managers:    managers,
		templateDir: templateDir,
		hasBindata:  hasBindata,
	}
}

type MultiManager struct {
	managers    []driver.Manager
	templateDir string
	hasBindata  bool
}

func (b *MultiManager) GetManagers() []driver.Manager {
	return b.managers
}

func (b *MultiManager) TemplateDir() string {
	return b.templateDir
}

func (b *MultiManager) HasBindata() bool {
	return b.hasBindata
}

func (b *MultiManager) Start() error {
	errs := common.NewErrors()
	for _, mgr := range b.managers {
		err := mgr.Start()
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}

func (b *MultiManager) Close() {
	for _, mgr := range b.managers {
		mgr.Close()
	}
}

func (b *MultiManager) ClearCallback() {
	for _, mgr := range b.managers {
		mgr.ClearCallback()
	}
}

func (b *MultiManager) AddCallback(rootDir string, callback func(name, typ, event string)) {
	originalCb := callback
	callback = func(name, typ, event string) {
		// 在 bindata 模式下，缓存的模板路径是类似于 “frontend/default/index.html” 的相对路径
		// 而 manager 监控的文件是绝对路径，所以需要裁剪成相对路径
		name = strings.TrimPrefix(name, b.templateDir)
		prefix := filepath.Base(strings.TrimSuffix(b.templateDir, echo.FilePathSeparator))
		name = path.Join(prefix, name)
		originalCb(name, typ, event)
	}
	for _, mgr := range b.managers {
		if !b.hasBindata {
			mgr.AddCallback(rootDir, originalCb)
			continue
		}
		if _, ok := mgr.(*manager.Manager); ok {
			mgr.AddCallback(rootDir, callback)
		} else {
			mgr.AddCallback(rootDir, originalCb)
		}
	}
}

func (b *MultiManager) DelCallback(rootDir string) {
	for _, mgr := range b.managers {
		mgr.DelCallback(rootDir)
	}
}

func (b *MultiManager) ClearAllows() {
	for _, mgr := range b.managers {
		mgr.ClearAllows()
	}
}

func (b *MultiManager) AddAllow(allows ...string) {
	for _, mgr := range b.managers {
		mgr.AddAllow(allows...)
	}
}

func (b *MultiManager) DelAllow(allow string) {
	for _, mgr := range b.managers {
		mgr.DelAllow(allow)
	}
}

func (b *MultiManager) ClearIgnores() {
	for _, mgr := range b.managers {
		mgr.ClearIgnores()
	}
}

func (b *MultiManager) AddIgnore(ignores ...string) {
	for _, mgr := range b.managers {
		mgr.AddIgnore(ignores...)
	}
}

func (b *MultiManager) DelIgnore(ignore string) {
	for _, mgr := range b.managers {
		mgr.DelIgnore(ignore)
	}
}

func (b *MultiManager) AddWatchDir(ppath string) (err error) {
	for _, mgr := range b.managers {
		err = mgr.AddWatchDir(ppath)
	}
	return
}

func (b *MultiManager) CancelWatchDir(oldDir string) (err error) {
	for _, mgr := range b.managers {
		err = mgr.CancelWatchDir(oldDir)
	}
	return
}

func (b *MultiManager) ChangeWatchDir(oldDir string, newDir string) (err error) {
	for _, mgr := range b.managers {
		err = mgr.ChangeWatchDir(oldDir, newDir)
	}
	return
}

func (b *MultiManager) SetLogger(logger logger.Logger) {
	for _, mgr := range b.managers {
		mgr.SetLogger(logger)
	}
}

func (b *MultiManager) ClearCache() {
	for _, mgr := range b.managers {
		mgr.ClearCache()
	}
}

func (b *MultiManager) GetTemplate(filename string) (content []byte, err error) {
	for _, mgr := range b.managers {
		filePath := filename
		if b.hasBindata {
			// 在 bindata 模式下，模板根目录是 template/ ， filename 会自动被追加上 frontend 或 backend 前缀
			// 原始模板的根目录是 template/frontend 或 template/backend ，所以在这个时候需要删除被追加的前缀
			if _, ok := mgr.(*manager.Manager); ok {
				filePath = filepath.Join(b.templateDir, strings.SplitN(filePath, `/`, 2)[1])
			}
		}
		content, err = mgr.GetTemplate(filePath)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}
	return
}

func (b *MultiManager) SetTemplate(filename string, content []byte) error {
	errs := common.NewErrors()
	for _, mgr := range b.managers {
		filePath := filename
		if b.hasBindata {
			// 在 bindata 模式下，模板根目录是 template/ ， filename 会自动被追加上 frontend 或 backend 前缀
			// 原始模板的根目录是 template/frontend 或 template/backend ，所以在这个时候需要删除被追加的前缀
			if _, ok := mgr.(*manager.Manager); ok {
				filePath = filepath.Join(b.templateDir, strings.SplitN(filePath, `/`, 2)[1])
			}
		}
		err := mgr.SetTemplate(filePath, content)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}
