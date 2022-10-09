package maker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
)

// MakeTemplate 生成模板
func MakeTemplate(cfg *Config) error {
	dir := TemplateDir
	if len(cfg.Group) > 0 {
		dir = filepath.Join(dir, cfg.Group)
	}
	err := com.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	tmplName := cfg.H.FileName()
	modelFile := filepath.Join(SampleDir, `template`)
	err = filepath.Walk(modelFile, func(ppath string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return err
		}
		name := fi.Name()
		if len(tmplName) > 0 {
			name = tmplName + `_` + fi.Name()
		}
		save := filepath.Join(dir, name)
		if com.FileExists(save) {
			fmt.Println(`Found:`, save, `Skipped`)
			return err
		}

		b, err := compile(ppath, cfg)
		if err != nil {
			return err
		}
		err = os.WriteFile(save, b, os.ModePerm)
		if err != nil {
			return err
		}
		err = Format(save)
		return err
	})
	return err
}
