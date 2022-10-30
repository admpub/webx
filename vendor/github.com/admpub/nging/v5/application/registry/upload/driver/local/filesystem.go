/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package local

import (
	"context"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	"github.com/admpub/nging/v5/application/registry/upload"
)

const Name = `local`

var _ upload.Storer = &Filesystem{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, subdir string) (upload.Storer, error) {
		return NewFilesystem(ctx, subdir), nil
	})
}

func NewFilesystem(ctx context.Context, subdir string, baseURLs ...string) *Filesystem {
	var baseURL string
	if len(baseURLs) > 0 {
		baseURL = baseURLs[0]
		baseURL = strings.TrimSuffix(baseURL, `/`)
	}
	return &Filesystem{
		Context: ctx,
		Subdir:  subdir,
		baseURL: baseURL,
	}
}

// Filesystem 文件系统存储引擎
type Filesystem struct {
	context.Context `json:"-" xml:"-"`
	Subdir          string
	baseURL         string
}

// Name 引擎名
func (f *Filesystem) Name() string {
	return Name
}

func (f *Filesystem) ErrIsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// FileDir 物理路径文件夹
func (f *Filesystem) FileDir(subpath string) string {
	return filepath.Join(uploadLibrary.UploadDir, f.Subdir, subpath)
}

// URLDir 网址路径文件夹
func (f *Filesystem) URLDir(subpath string) string {
	return path.Join(uploadLibrary.UploadURLPath, f.Subdir, subpath)
}

// Exists 判断文件是否存在
func (f *Filesystem) Exists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileInfo 获取文件信息
func (f *Filesystem) FileInfo(file string) (os.FileInfo, error) {
	return os.Stat(file)
}

// SendFile 下载文件
func (f *Filesystem) SendFile(ctx echo.Context, file string) error {
	return ctx.File(file)
}

// Put 上传文件
func (f *Filesystem) Put(dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = f.FileDir(dstFile)
	saveDir := filepath.Dir(savePath)
	err = com.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return
	}
	viewURL = f.PublicURL(dstFile)
	//create destination file making sure the path is writeable.
	var dst *os.File
	dst, err = os.Create(savePath)
	if err != nil {
		return
	}
	defer dst.Close()
	//copy the uploaded file to the destination file
	_, err = io.Copy(dst, src)
	return
}

// PublicURL 文件物理路径转为文件网址
func (f *Filesystem) PublicURL(dstFile string) string {
	return f.baseURL + f.URLDir(dstFile)
}

// URLToFile 文件网址转为文件物理路径
func (f *Filesystem) URLToFile(publicURL string) string {
	dstFile := f.URLToPath(publicURL)
	dstFile = strings.TrimPrefix(dstFile, strings.TrimRight(f.URLDir(``), `/`)+`/`)
	return dstFile
}

// URLToPath 文件网址转为文件路径
func (f *Filesystem) URLToPath(publicURL string) string {
	if len(f.baseURL) > 0 {
		publicURL = strings.TrimPrefix(publicURL, f.baseURL+`/`)
		if !strings.HasPrefix(publicURL, `/`) {
			publicURL = `/` + publicURL
		}
	}
	return publicURL
}

// SetBaseURL 设置根网址
func (f *Filesystem) SetBaseURL(baseURL string) {
	baseURL = strings.TrimSuffix(baseURL, `/`)
	f.baseURL = baseURL
}

// BaseURL 根网址
func (f *Filesystem) BaseURL() string {
	return f.baseURL
}

// FixURL 改写文件网址
func (f *Filesystem) FixURL(content string, embedded ...bool) string {
	return content
}

// URLWithParams 文件网址增加参数
func (f *Filesystem) URLWithParams(rawURL string, values url.Values) string {
	if values == nil {
		return rawURL
	}
	queryString := values.Encode()
	if len(queryString) > 0 {
		rawURL += `?`
	}
	rawURL += queryString
	return rawURL
}

// Get 获取文件读取接口
func (f *Filesystem) Get(dstFile string) (io.ReadCloser, error) {
	return f.openFile(dstFile)
}

func (f *Filesystem) openFile(dstFile string) (*os.File, error) {
	//file := f.filepath(dstFile)
	file := filepath.Join(echo.Wd(), dstFile)
	return os.Open(file)
}

// Delete 删除文件
func (f *Filesystem) Delete(dstFile string) error {
	file := filepath.Join(echo.Wd(), dstFile)
	return os.Remove(file)
}

// DeleteDir 删除文件夹及其内部文件
func (f *Filesystem) DeleteDir(dstDir string) error {
	dir := filepath.Join(echo.Wd(), dstDir)
	return os.RemoveAll(dir)
}

// Move 移动文件
func (f *Filesystem) Move(src, dst string) error {
	sdir := filepath.Dir(dst)
	if err := com.MkdirAll(sdir, os.ModePerm); err != nil {
		return err
	}
	return com.Rename(src, dst)
}

// Close 关闭连接
func (f *Filesystem) Close() error {
	return nil
}
