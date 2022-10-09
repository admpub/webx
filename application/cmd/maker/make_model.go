package maker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
)

// MakeModel 生成模型
func MakeModel(cfg *Config) error {
	dir := ModelDir
	if len(cfg.Group) > 0 {
		dir = filepath.Join(dir, cfg.Group)
	}
	err := com.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	name := cfg.M.FileName()
	save := filepath.Join(dir, name+`.go`)

	if com.FileExists(save) {
		fmt.Println(`Found:`, save, `Skipped`)
		return nil
	}
	modelFile := filepath.Join(SampleDir, `model`, `model.go.tpl`)
	b, err := compile(modelFile, cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(save, b, os.ModePerm)
	if err != nil {
		return err
	}
	err = Format(save)
	return err
}
