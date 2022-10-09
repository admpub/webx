package iotools

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SafeFile struct {
	*os.File
	lock     sync.Mutex
	filePath string
	closed   bool
}

func (sf *SafeFile) FilePath() string {
	return sf.filePath
}

func (sf *SafeFile) SetFilePath(filePath string) *SafeFile {
	sf.filePath = filePath
	return sf
}

func (sf *SafeFile) WriteAt(b []byte, off int64) (n int, err error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	return sf.File.WriteAt(b, off)
}

func (sf *SafeFile) Sync() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	return sf.File.Sync()
}

func (sf *SafeFile) Close() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	if sf.closed {
		return nil
	}
	err := sf.File.Close()
	if err == nil {
		sf.closed = true
	} else {
		log.Println(`-> close file`, sf.filePath, `error:`, err)
	}
	return err
}

func (sf *SafeFile) ReOpen() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	if !sf.closed {
		return nil
	}
	f, err := os.OpenFile(sf.filePath, os.O_RDWR, 0666)
	if err == nil {
		sf.closed = false
	} else if os.IsNotExist(err) {
		var errc error
		f, errc = os.Create(sf.filePath)
		if errc != nil {
			log.Println(`=> create file`, sf.filePath, `error:`, errc)
		}
		sf.closed = false
	} else {
		log.Println(`=> open file`, sf.filePath, `error:`, err)
	}
	sf.File = f
	return err
}

func OpenSafeFile(name string) (file *SafeFile, err error) {
	file = &SafeFile{filePath: name, closed: true}
	err = file.ReOpen()
	return
}

func CreateSafeFile(name string) (file *SafeFile, err error) {
	var f *os.File
	name = strings.TrimRight(name, `/`)
	name = strings.TrimRight(name, `\`)
	dir := filepath.Dir(name)
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		os.MkdirAll(dir, 0666)
	}
	f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	file = &SafeFile{File: f, filePath: name}
	return
}
