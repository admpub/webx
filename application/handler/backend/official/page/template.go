package page

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/library/xtemplate"
)

func TemplatePreviewImage(ctx echo.Context) error {
	name := ctx.Query(`name`)
	color := ctx.Query(`color`)
	themeInfo, err := getTemplateInfo(name)
	if err != nil {
		return err
	}
	previewImage := themeInfo.PreviewImage
	if len(color) > 0 {
		previewImage = ``
		for _, tc := range themeInfo.Colors {
			if tc.Name == color {
				previewImage = tc.PreviewImage
				break
			}
		}
	}
	if len(previewImage) == 0 {
		return echo.ErrNotFound
	}
	previewImageFile := filepath.Join(name, previewImage)
	file, err := GetTemplateDiskFS().Open(previewImageFile)
	if err == nil {
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		return ctx.ServeContent(file, previewImage, fi.ModTime())
	}
	if GetTemplateEmbedFS() != nil {
		previewImageFile = path.Join(templateRoot, name, previewImage)

		if tfs, err := GetTemplateEmbedFS().Open(previewImageFile); err == nil {
			defer tfs.Close()
			fi, err := tfs.Stat()
			if err != nil {
				return err
			}
			return ctx.ServeContent(tfs, previewImage, fi.ModTime())
		}
	}
	return echo.ErrNotFound
}

func TemplateEnable(ctx echo.Context) error {
	name := ctx.Form(`name`)
	themeInfo, err := getTemplateInfo(name)
	data := ctx.Data()
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	var cfg *xtemplate.ThemeInfo
	cfg, err = frontend.TmplPathFixers.Storer().Get(ctx, name)
	if err != nil {
		if !os.IsNotExist(err) {
			return ctx.JSON(data.SetError(err))
		}
	} else {
		themeInfo.CustomConfig.DeepMerge(cfg.CustomConfig)
	}
	frontend.TmplPathFixers.SetThemeInfo(ctx, themeInfo)
	return ctx.JSON(data)
}

var (
	themeList       []*xtemplate.ThemeInfo
	themeLsMu       sync.RWMutex
	themeLsInitOnce sync.Once
)

func getTemplateList() []*xtemplate.ThemeInfo {
	var (
		dirs   []fs.FileInfo
		list   []*xtemplate.ThemeInfo
		pdirs  = map[string]struct{}{}
		embeds []*xtemplate.ThemeInfo
	)
	dirs, err := GetTemplateDiskFS().ReadDir(`./`)
	if err != nil && !os.IsNotExist(err) {
		return list
	}
	list = make([]*xtemplate.ThemeInfo, 0, len(dirs))
	embeds = GetEmbedThemes()
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), `.`) || !dir.IsDir() {
			continue
		}
		themeInfo := &xtemplate.ThemeInfo{
			Name: dir.Name(),
		}
		infoFile := filepath.Join(dir.Name(), `@info.yaml`)
		content, err := GetTemplateDiskFS().ReadFile(infoFile)
		if err == nil {
			themeInfo.Decode(content)
		} else {
			for _, v := range embeds {
				if v.Name == themeInfo.Name {
					themeInfo = v
					break
				}
			}
			if !themeInfo.Embed() {
				themeInfo.Author.Name = `coscms`
				themeInfo.Title = dir.Name()
				themeInfo.Version = `0.0.1`
				themeInfo.UpdatedAt = time.Now().Format(param.DateTimeNormal)
				themeInfo.Fallback = []string{`default`}
				infoFile = filepath.Join(frontend.DefaultTemplateDir, infoFile)
				themeInfo.EncodeToFile(infoFile)
			}
		}
		list = append(list, themeInfo)
		pdirs[themeInfo.Name] = struct{}{}
	}
	for _, v := range embeds {
		if _, ok := pdirs[v.Name]; ok {
			continue
		}
		list = append(list, v)
	}
	return list
}

func TemplateIndex(ctx echo.Context) error {
	if ctx.Form(`op`) == `preview` {
		return TemplatePreviewImage(ctx)
	}
	list := getTemplateList()
	themeLsMu.Lock()
	themeList = list
	themeLsInitOnce.Do(func() {})
	themeLsMu.Unlock()

	ctx.Set(`listData`, list)
	ctx.Set(`current`, frontend.TmplPathFixers.ThemeInfo(ctx))
	return ctx.Render(`official/page/template_index`, nil)
}

var canEditExtensions = []string{`.js`, `.html`, `.css`, `.json`, `.txt`}
var fileNameRegex = regexp.MustCompile(`^([^/\"']+)(\.[a-zA-Z0-9]+)?|\.[a-zA-Z0-9]+$`)

func valiateFileName(ctx echo.Context, fileName string) error {
	if fileNameRegex.MatchString(fileName) {
		return nil
	}
	return ctx.NewError(code.InvalidParameter, `文件名格式不正确。文件名不能包含引号和正、反斜杠，扩展名只能包含字母和数字`)
}

func TemplateEdit(ctx echo.Context) error {
	name := ctx.Form(`name`)
	if len(name) == 0 {
		return echo.ErrNotFound
	}
	if !xtemplate.IsThemeName(name) {
		return echo.ErrNotFound
	}
	themeDir := name
	var dirPositions []string
	var dirURLs []string
	var embedThemeDir string
	dir := ctx.Form(`dir`)
	if len(dir) > 0 {
		dir = filepath.Clean(dir)
		themeDir = filepath.Join(themeDir, dir)
	}
	var tfs http.File
	var err error

	fi, err := GetTemplateDiskFS().Stat(themeDir)
	exists := err == nil && fi.IsDir()
	closeEmbedFS := func() error {
		if tfs != nil {
			return tfs.Close()
		}
		return nil
	}
	defer closeEmbedFS()
	getEmbedFS := func() http.File {
		if tfs != nil {
			return tfs
		}
		if GetTemplateEmbedFS() == nil {
			return nil
		}
		embedThemeDir = path.Join(templateRoot, name)
		if len(dir) > 0 {
			embedThemeDir = path.Join(embedThemeDir, dir)
		}
		tfs, err = GetTemplateEmbedFS().Open(embedThemeDir)
		if err != nil {
			err = fmt.Errorf(`%s: %w`, embedThemeDir, err)
		}
		return tfs
	}
	if !exists && GetTemplateEmbedFS() != nil {
		if getEmbedFS() == nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			exists = true
		}
	}
	if !exists {
		log.Debugf(`%v: %s`, echo.ErrNotFound, themeDir)
		return echo.ErrNotFound
	}
	op := ctx.Form(`op`)
	switch op {
	case `getFileContent`:
		file := ctx.Form(`file`)
		if len(file) == 0 {
			return ctx.NewError(code.InvalidParameter, `请选择一个要编辑的文件`).SetZone(`file`)
		}
		_file := file
		var b []byte
		file = filepath.Join(themeDir, file)
		b, err := GetTemplateDiskFS().ReadFile(file)
		if err == nil {
			return ctx.JSON(ctx.Data().SetData(echo.H{`content`: com.Bytes2str(b)}))
		}
		if getEmbedFS() == nil {
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			return echo.ErrNotFound
		}
		file = path.Clean(_file)
		file = path.Join(embedThemeDir, file)
		var f http.File
		f, err = GetTemplateEmbedFS().Open(file)
		if err != nil {
			return err
		}
		b, err = io.ReadAll(f)
		f.Close()
		if err != nil {
			return err
		}
		return ctx.JSON(ctx.Data().SetData(echo.H{`content`: com.Bytes2str(b)}))
	case `saveFileContent`: // 保存文件内容
		file := ctx.Form(`file`)
		if len(file) == 0 {
			return ctx.NewError(code.InvalidParameter, `请选择一个要保存的文件`).SetZone(`file`)
		}
		file = filepath.Clean(file)
		if err := valiateFileName(ctx, file); err != nil {
			return err
		}
		original := file
		themeDir = filepath.Join(frontend.DefaultTemplateDir, themeDir)
		file = filepath.Join(themeDir, file)
		if ctx.Formx(`isNew`).Bool() {
			if com.FileExists(file) && !ctx.Formx(`confirmed`).Bool() {
				return ctx.NewError(code.DataAlreadyExists, `文件“%s”已经存在，确定要覆盖吗？`, original)
			}
		} else if !com.FileExists(file) {
			if tfs == nil {
				log.Debugf(`%v: %s`, echo.ErrNotFound, file)
				return echo.ErrNotFound
			}
			com.MkdirAll(themeDir, os.ModePerm)
		}
		content := ctx.Form(`content`)
		err := os.WriteFile(file, com.Str2bytes(content), os.ModePerm)
		if err != nil {
			return err
		}
		err = echo.FireByName(EventNameFrontendTemplateEdited, events.WithContext(echo.H{`file`: file, `op`: op}))
		if err != nil {
			return err
		}
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`保存成功`), 1))
	case `saveAsFileContent`: // 将文件内容另存为新文件
		file := ctx.Form(`file`)
		if len(file) == 0 {
			return ctx.NewError(code.InvalidParameter, `请设置一个要另存为的新文件名称`).SetZone(`file`)
		}
		file = filepath.Clean(file)
		if err := valiateFileName(ctx, file); err != nil {
			return err
		}
		original := file
		themeDir = filepath.Join(frontend.DefaultTemplateDir, themeDir)
		file = filepath.Join(themeDir, file)
		if com.FileExists(file) && !ctx.Formx(`confirmed`).Bool() {
			return ctx.NewError(code.DataAlreadyExists, `文件“%s”已经存在，确定要覆盖吗？`, original)
		}
		if tfs != nil {
			com.MkdirAll(themeDir, os.ModePerm)
		}
		content := ctx.Form(`content`)
		err := os.WriteFile(file, com.Str2bytes(content), os.ModePerm)
		if err != nil {
			return err
		}
		err = echo.FireByName(EventNameFrontendTemplateEdited, events.WithContext(echo.H{`file`: file, `op`: op}))
		if err != nil {
			return err
		}
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`保存成功`), 1))
	case `renameFile`:
		file := ctx.Form(`file`)
		if len(file) == 0 {
			return ctx.NewError(code.InvalidParameter, `请选择一个要改名的文件`).SetZone(`file`)
		}
		newFile := ctx.Form(`newFile`)
		if len(file) == 0 {
			return ctx.NewError(code.InvalidParameter, `请设置文件的新名称`).SetZone(`file`)
		}
		original := newFile
		srcFile := filepath.Join(themeDir, file)
		themeDir = filepath.Join(frontend.DefaultTemplateDir, themeDir)
		newFile = filepath.Join(themeDir, newFile)
		if com.FileExists(newFile) && !ctx.Formx(`confirmed`).Bool() {
			return ctx.NewError(code.DataAlreadyExists, `文件“%s”已经存在，确定要覆盖吗？`, original)
		}
		newFile = filepath.Clean(newFile)
		if err = valiateFileName(ctx, newFile); err != nil {
			return err
		}
		var b []byte
		if getEmbedFS() == nil {
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			file = filepath.Clean(file)
			if file == newFile {
				return ctx.JSON(ctx.Data().SetInfo(ctx.T(`保存成功`), 1))
			}
			b, err = GetTemplateDiskFS().ReadFile(srcFile)
		} else {
			file = path.Clean(file)
			if file == newFile {
				return ctx.JSON(ctx.Data().SetInfo(ctx.T(`保存成功`), 1))
			}
			file = path.Join(embedThemeDir, file)
			var f http.File
			f, err = GetTemplateEmbedFS().Open(file)
			if err != nil {
				if os.IsNotExist(err) {
					return ctx.JSON(ctx.Data().SetInfo(ctx.T(`内置文件“%s”不存在`, file), 0).SetZone(`file`))
				}
				return err
			}
			b, err = io.ReadAll(f)
			f.Close()
		}
		if err == nil {
			err = os.WriteFile(newFile, b, os.ModePerm)
		}
		if err != nil {
			return err
		}
		err = echo.FireByName(EventNameFrontendTemplateEdited, events.WithContext(echo.H{`file`: file, `newFile`: newFile, `op`: op}))
		if err != nil {
			return err
		}
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`保存成功`), 1))
	case `removeFile`:
		if tfs != nil {
			return ctx.NewError(code.InvalidParameter, `不能删除内置模板文件`).SetZone(`file`)
		}
		file := ctx.Form(`file`)
		if len(file) == 0 {
			return ctx.NewError(code.InvalidParameter, `请选择一个要删除的文件`).SetZone(`file`)
		}
		file = filepath.Clean(file)
		original := file
		file = filepath.Join(frontend.DefaultTemplateDir, themeDir, file)
		if !com.FileExists(file) {
			return ctx.NewError(code.DataAlreadyExists, `文件“%s”不存在或不属于当前项目`, original)
		}
		err := os.Remove(file)
		if err != nil {
			return err
		}
		err = echo.FireByName(EventNameFrontendTemplateEdited, events.WithContext(echo.H{`file`: file, `op`: op}))
		if err != nil {
			return err
		}
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`删除成功`), 1))
	}
	var (
		dirs  []fs.FileInfo
		ldirs []xtemplate.FileInfo
		pdirs = map[string]struct{}{}
	)
	if tfs == nil && GetTemplateEmbedFS() != nil {
		embedThemeDir = path.Join(templateRoot, name)
		if len(dir) > 0 {
			embedThemeDir = path.Join(templateRoot, dir)
		}
		tfs, err = GetTemplateEmbedFS().Open(embedThemeDir)
		if err == nil {
			defer tfs.Close()
		}
	}
	dirs, err = GetTemplateDiskFS().ReadDir(themeDir)
	if err != nil {
		if tfs != nil {
			goto EMD
		}
		goto END
	}
	for _, d := range dirs {
		ldirs = append(ldirs, xtemplate.FileInfo{FileInfo: d})
		pdirs[d.Name()] = struct{}{}
	}

EMD:
	if getEmbedFS() != nil {
		var fis []fs.FileInfo
		fis, err = tfs.Readdir(-1)
		if err != nil {
			goto END
		}
		for _, v := range fis {
			if _, ok := pdirs[v.Name()]; ok {
				continue
			}
			ldirs = append(ldirs, xtemplate.FileInfo{FileInfo: v, Embed: true})
		}
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	sort.Sort(xtemplate.SortFileInfoByFileType(ldirs))
	ctx.Set(`dirs`, ldirs)
	if len(dir) > 0 {
		dirPositions = strings.Split(strings.Trim(filepath.ToSlash(dir), `/`), `/`)
		dirURLs = make([]string, len(dirPositions))
		curURL := handler.URLFor(`/official/page/template_edit`) + `?name=` + url.QueryEscape(name)
		for index, dirName := range dirPositions {
			var prefix string
			if index > 0 {
				prefix = strings.Join(dirPositions[0:index], `/`) + `/`
			}
			dirURLs[index] = curURL + `&dir=` + url.QueryEscape(prefix+dirName)
		}
	}

END:
	ctx.Set(`activeURL`, `/official/page/template_index`)
	ctx.Set(`dirPositions`, dirPositions)
	ctx.Set(`dirURLs`, dirURLs)
	ctx.SetFunc(`canEdit`, func(file string) bool {
		return com.InSlice(filepath.Ext(file), canEditExtensions)
	})
	return ctx.Render(`official/page/template_edit`, handler.Err(ctx, err))
}

func TemplateConfig(ctx echo.Context) error {
	name := ctx.Form(`name`)
	current := frontend.TmplPathFixers.ThemeInfo(ctx)
	var (
		cfg          *xtemplate.ThemeInfo
		defaultColor string
	)
	themeInfo, err := getTemplateInfo(name)
	if err != nil {
		goto END
	}
	cfg, err = frontend.TmplPathFixers.Storer().Get(ctx, name)
	if err != nil {
		if !os.IsNotExist(err) {
			goto END
		}
		err = nil
	} else {
		themeInfo.CustomConfig.DeepMerge(cfg.CustomConfig)
		if len(cfg.Fallback) > 0 {
			themeInfo.Fallback = cfg.Fallback
		}
	}
	if ctx.IsPost() {
		themeInfo.CustomConfig.Set(`color`, ctx.Form(`color`))
		err = themeInfo.SaveForm(ctx, `config`)
		if err != nil {
			goto END
		}
		fallbackThemes := ctx.Form(`_fallbackThemes`)
		if len(fallbackThemes) > 0 {
			themeInfo.Fallback = param.StringSlice(strings.Split(fallbackThemes, `,`)).Filter(func(s *string) bool {
				if s == nil || len(*s) == 0 || *s == themeInfo.Name || *s == `default` {
					return false
				}
				return true
			}).Unique().String()
		} else {
			themeInfo.Fallback = []string{}
		}
		if themeInfo.Name != `default` {
			themeInfo.Fallback = append(themeInfo.Fallback, `default`)
		}
		err = frontend.TmplPathFixers.Storer().Put(ctx, name, themeInfo)
		if err != nil {
			goto END
		}
		if current.Name == name {
			frontend.TmplPathFixers.SetThemeInfo(ctx, themeInfo)
		}
		handler.SendOk(ctx, ctx.T(`保存成功`))
		return ctx.Redirect(handler.URLFor(`/official/page/template_index`))
	}
	defaultColor = themeInfo.CustomConfig.String(`color`)
	if len(defaultColor) == 0 && len(themeInfo.Colors) > 0 {
		for _, v := range themeInfo.Colors {
			if v.IsDefault {
				defaultColor = v.Name
				break
			}
		}
		if len(defaultColor) == 0 {
			defaultColor = themeInfo.Colors[0].Name
		}
		if len(defaultColor) > 0 {
			themeInfo.CustomConfig.Set(`color`, defaultColor)
		}
	}
	for k, v := range themeInfo.CustomConfig {
		ctx.Request().Form().Set(k, param.AsString(v))
	}
END:
	ctx.Set(`info`, themeInfo)
	ctx.Set(`activeURL`, `/official/page/template_index`)

	themeLsInitOnce.Do(func() {
		themeLsMu.Lock()
		themeList = getTemplateList()
		themeLsMu.Unlock()
	})

	themeLsMu.RLock()
	fallbacks := make([]xtemplate.ThemeInfoLite, 0, len(themeList))
	for _, themeCfg := range themeList {
		if themeCfg.Name == themeInfo.Name || themeCfg.Name == `default` {
			continue
		}
		lite := themeCfg.AsLite()
		if len(lite.PreviewImage) > 0 {
			lite.PreviewImage = handler.URLFor(`/official/page/template_index`) + `?op=preview&name=` + lite.Name
		}
		fallbacks = append(fallbacks, lite)
	}
	themeLsMu.RUnlock()

	ctx.Set(`fallbacks`, fallbacks)
	return ctx.Render(`official/page/template_config`, handler.Err(ctx, err))
}
