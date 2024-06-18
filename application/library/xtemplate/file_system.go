package xtemplate

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var _ http.FileSystem = NewFileSystems()

func NewFileSystems() FileSystems {
	return FileSystems{}
}

type FileSystems []http.FileSystem

func (f FileSystems) Open(name string) (file http.File, err error) {
	for _, fileSystem := range f {
		file, err = fileSystem.Open(name)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}
	return
}

func (f FileSystems) ReadFile(name string) (content []byte, err error) {
	var fp http.File
	fp, err = f.Open(name)
	if err != nil {
		return
	}
	if fp == nil {
		err = os.ErrNotExist
		return
	}
	b, err := io.ReadAll(fp)
	fp.Close()
	return b, err
}

func (f FileSystems) ReadDir(name string) (dirs []fs.FileInfo, err error) {
	unique := map[string]struct{}{}
	for _, fileSystem := range f {
		var file http.File
		file, err = fileSystem.Open(name)
		if err != nil {
			if !os.IsNotExist(err) {
				return
			}
			err = nil
			continue
		}
		if fi, err := file.Stat(); err != nil || !fi.IsDir() {
			file.Close()
			continue
		}
		var _dirs []fs.FileInfo
		_dirs, err = file.Readdir(-1)
		file.Close()
		if err != nil {
			if !os.IsNotExist(err) {
				return
			}
			err = nil
			continue
		}
		for _, dir := range _dirs {
			if _, ok := unique[dir.Name()]; ok {
				continue
			}
			unique[dir.Name()] = struct{}{}
			dirs = append(dirs, dir)
		}
	}
	return
}

func (f FileSystems) Stat(name string) (fi fs.FileInfo, err error) {
	var fp http.File
	fp, err = f.Open(name)
	if err != nil {
		return
	}
	if fp == nil {
		err = os.ErrNotExist
		return
	}
	fi, err = fp.Stat()
	fp.Close()
	return
}

func (f FileSystems) Size() int {
	return len(f)
}

func (f FileSystems) IsEmpty() bool {
	return f.Size() == 0
}

func (f *FileSystems) Register(fileSystem http.FileSystem) {
	*f = append(*f, fileSystem)
}

func NewFileSystemTrimPrefix(trimPrefix string, fs http.FileSystem) http.FileSystem {
	if len(trimPrefix) == 0 {
		return fs
	}
	return fsTrimPrefix{
		trimPrefix: trimPrefix,
		fs:         fs,
	}
}

type fsTrimPrefix struct {
	trimPrefix string
	fs         http.FileSystem
}

func (f fsTrimPrefix) Open(name string) (file http.File, err error) {
	name = strings.TrimPrefix(name, f.trimPrefix)
	return f.fs.Open(name)
}

func NewStaticDir(dir string, trimPrefix ...string) http.FileSystem {
	var _trimPrefix string
	if len(trimPrefix) > 0 {
		_trimPrefix = trimPrefix[0]
	}
	return &StaticDir{
		dir:        dir,
		trimPrefix: _trimPrefix,
	}
}

type StaticDir struct {
	dir        string
	trimPrefix string
}

func (s *StaticDir) Open(name string) (http.File, error) {
	if len(s.trimPrefix) > 0 {
		name = strings.TrimPrefix(name, s.trimPrefix)
	}
	fullName := filepath.Join(s.dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	return f, err
}

type FileInfo struct {
	fs.FileInfo
	Embed bool
}

type SortFileInfoByFileType []FileInfo

func (s SortFileInfoByFileType) Len() int { return len(s) }
func (s SortFileInfoByFileType) Less(i, j int) bool {
	if s[i].IsDir() {
		if !s[j].IsDir() {
			return true
		}
	} else if s[j].IsDir() {
		if !s[i].IsDir() {
			return false
		}
	}
	return s[i].Name() < s[j].Name()
}
func (s SortFileInfoByFileType) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
