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

func (f FileSystems) ReadDir(name string, count int) (dirs []fs.FileInfo, err error) {
	var fp http.File
	fp, err = f.Open(name)
	if err != nil {
		return
	}
	if fp == nil {
		err = os.ErrNotExist
		return
	}
	dirs, err = fp.Readdir(count)
	fp.Close()
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

type DirEntry struct {
	fs.DirEntry
	Embed bool
}

type SortDirEntryByFileType []DirEntry

func (s SortDirEntryByFileType) Len() int { return len(s) }
func (s SortDirEntryByFileType) Less(i, j int) bool {
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
func (s SortDirEntryByFileType) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
