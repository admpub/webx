package page

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/events"
	"github.com/admpub/log"
	formConfig "github.com/coscms/forms/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/bindata"
)

const EventNameFrontendTemplateEdited = `frontend.template.edited`

func init() {
	echo.OnCallback(EventNameFrontendTemplateEdited, func(e events.Event) error {
		httpserver.Frontend.Template.ClearCache()
		return echo.FireByName(`webx.renderer.cache.clear`)
	})
}

func getTemplateInfo(name string) (*ntemplate.ThemeInfo, error) {
	if len(name) == 0 {
		return nil, echo.ErrNotFound
	}
	if !ntemplate.IsThemeName(name) {
		return nil, echo.ErrNotFound
	}
	themeInfo := &ntemplate.ThemeInfo{
		CustomConfig: echo.H{},
		FormConfig:   make(map[string]formConfig.Config),
	}
	infoFile := filepath.Join(name, `@info.yaml`)
	content, err := GetTemplateDiskFS().ReadFile(infoFile)
	if err == nil {
		err = themeInfo.Decode(content)
		return themeInfo, err
	}
	if GetTemplateEmbedFS() != nil {
		infoFile = path.Join(templateRoot, name, `@info.yaml`)

		if tfs, err := GetTemplateEmbedFS().Open(infoFile); err == nil {
			b, err := io.ReadAll(tfs)
			tfs.Close()
			if err == nil {
				err = themeInfo.Decode(b)
				return themeInfo, err
			}
		}
	}
	return nil, echo.ErrNotFound
}

var (
	templateRoot    = `/frontend/`
	templateDiskFS  ntemplate.FileSystems
	templateDiskMx  sync.Once
	templateEmbedFS http.FileSystem
	templateEmbedMx sync.Once
	embedThemes     []*ntemplate.ThemeInfo
	embedThemesMx   sync.Once
)

func initTemplateDiskFS() {
	templateDiskFS = ntemplate.NewFileSystems()
	templateDiskFS.Register(http.Dir(frontend.DefaultTemplateDir))
	initTemplateDiskOtherFS()
}

func initTemplateEmbedFS() {
	switch mgr := httpserver.Frontend.TmplMgr.(type) {
	case *bindata.TmplManager:
		templateEmbedFS = mgr.FileSystem
	case *ntemplate.MultiManager:
		for _, m := range mgr.GetManagers() {
			if mgr, ok := m.(*bindata.TmplManager); ok {
				templateEmbedFS = mgr.FileSystem
				break
			}
		}
	}
}

func initEmbedThemes() {
	if GetTemplateEmbedFS() == nil {
		return
	}

	tfs, err := GetTemplateEmbedFS().Open(templateRoot)
	if err != nil {
		return
	}
	defer tfs.Close()
	dirs, _ := tfs.Readdir(-1)
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), `.`) || !dir.IsDir() {
			continue
		}
		infoFile := path.Join(templateRoot, dir.Name(), `@info.yaml`)
		file, err := GetTemplateEmbedFS().Open(infoFile)
		if err != nil {
			continue
		}
		b, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			continue
		}
		themeInfo := &ntemplate.ThemeInfo{
			Name: dir.Name(),
		}
		themeInfo.Decode(b)
		themeInfo.SetEmbed()
		embedThemes = append(embedThemes, themeInfo)
	}
}

func GetTemplateEmbedFS() http.FileSystem {
	templateEmbedMx.Do(initTemplateEmbedFS)
	return templateEmbedFS
}

func GetTemplateDiskFS() ntemplate.FileSystems {
	templateDiskMx.Do(initTemplateDiskFS)
	return templateDiskFS
}

func GetEmbedThemes() []*ntemplate.ThemeInfo {
	embedThemesMx.Do(initEmbedThemes)
	return embedThemes
}

func GetAllThemeNames() []string {
	list := getTemplateList()
	themes := make([]string, len(list))
	for i, v := range list {
		themes[i] = v.Name
	}
	return themes
}

func ListTemplateFiles(dir string, themes ...string) (r []fs.FileInfo) {
	if len(themes) == 0 {
		themes = GetAllThemeNames()
	}
	unique := map[string]struct{}{}
	embedFS := GetTemplateEmbedFS()
	for _, theme := range themes {
		full := filepath.Join(theme, dir)
		files, err := GetTemplateDiskFS().ReadDir(full)
		if err != nil {
			log.Error(err)
		}
		if embedFS != nil {
			embedFull := path.Join(templateRoot, theme, dir)
			embedFile, err := embedFS.Open(embedFull)
			if err != nil {
				log.Error(err)
			} else {
				dirs, _ := embedFile.Readdir(-1)
				for _, fInfo := range dirs {
					if fInfo.IsDir() {
						continue
					}
					files = append(files, fInfo)
				}
			}
		}
		for _, fi := range files {
			if _, ok := unique[fi.Name()]; ok {
				continue
			}
			unique[fi.Name()] = struct{}{}
			r = append(r, fi)
		}
	}
	return
}

func ListTemplateFileNames(dir string, themes ...string) []string {
	files := ListTemplateFiles(dir, themes...)
	if len(files) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(files))
	for _, info := range files {
		if info.IsDir() {
			continue
		}
		names = append(names, info.Name())
	}
	return names
}
