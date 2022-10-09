package maker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// MakeHandler 生成Handler
func MakeHandler(cfg *Config) error {
	dir := HandlerDir
	if len(cfg.Group) > 0 {
		dir = filepath.Join(dir, cfg.Group)
	}
	err := com.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	name := cfg.H.FileName()
	save := filepath.Join(dir, name+`.go`)
	if com.FileExists(save) {
		fmt.Println(`Found:`, save, `Skipped`)
		return nil
	}

	handlerFile := filepath.Join(SampleDir, `handler`, `handler.go.tpl`)
	b, err := compile(handlerFile, cfg)
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

// MakeHandlerInit 生成Handler初始化逻辑
func MakeHandlerInit(group string, data echo.H) error {
	dir := HandlerDir
	if len(group) > 0 {
		dir = filepath.Join(dir, group)
	}
	err := com.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	initFiles := []string{`init.go`, `init_navigate.go`}
	for _, initFile := range initFiles {
		save := filepath.Join(dir, initFile)
		if com.FileExists(save) {
			fmt.Println(`Found:`, save, `Skipped`)
			continue
		}

		handlerFile := filepath.Join(SampleDir, `handler`, initFile+`.tpl`)
		b, err := compile(handlerFile, data)
		if err != nil {
			return err
		}
		err = os.WriteFile(save, b, os.ModePerm)
		if err != nil {
			return err
		}
		err = Format(save)
	}
	return err
}
