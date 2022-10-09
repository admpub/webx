package page

import (
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/events"
	"github.com/admpub/nging/v4/application/cmd/bootconfig"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/xtemplate"
	formConfig "github.com/coscms/forms/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/bindata"
)

const EventNameFrontendTemplateEdited = `frontend.template.edited`

func init() {
	echo.OnCallback(EventNameFrontendTemplateEdited, func(e events.Event) error {
		frontend.TmplPathFixers.ClearCache()
		return echo.FireByName(`webx.renderer.cache.clear`)
	})
}

func getTemplateInfo(name string) (*xtemplate.ThemeInfo, error) {
	templateDir := frontend.TemplateDir
	if len(name) == 0 {
		return nil, echo.ErrNotFound
	}
	if !xtemplate.IsThemeName(name) {
		return nil, echo.ErrNotFound
	}
	themeInfo := &xtemplate.ThemeInfo{
		CustomConfig: echo.H{},
		FormConfig:   make(map[string]formConfig.Config),
	}
	infoFile := filepath.Join(templateDir, name, `@info.yaml`)
	if !com.FileExists(infoFile) {
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
	err := themeInfo.DecodeFile(infoFile)
	return themeInfo, err
}

var (
	templateRoot    = `/frontend/`
	templateEmbedFS http.FileSystem
	templateEmbedMx sync.Once
	embedThemes     []*xtemplate.ThemeInfo
	embedThemesMx   sync.Once
)

func initTemplateFS() {
	if mgr, ok := bootconfig.FrontendTmplMgr.(*bindata.TmplManager); ok {
		templateEmbedFS = mgr.FileSystem
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
		themeInfo := &xtemplate.ThemeInfo{
			Name: dir.Name(),
		}
		themeInfo.Decode(b)
		themeInfo.SetEmbed()
		embedThemes = append(embedThemes, themeInfo)
	}
}

func GetTemplateEmbedFS() http.FileSystem {
	templateEmbedMx.Do(initTemplateFS)
	return templateEmbedFS
}

func GetEmbedThemes() []*xtemplate.ThemeInfo {
	embedThemesMx.Do(initEmbedThemes)
	return embedThemes
}
