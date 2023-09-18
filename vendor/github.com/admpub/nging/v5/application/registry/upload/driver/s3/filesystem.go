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

package s3

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	"github.com/admpub/nging/v5/application/library/s3manager"
	"github.com/admpub/nging/v5/application/library/s3manager/s3client"
	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/nging/v5/application/registry/upload"
	"github.com/admpub/nging/v5/application/registry/upload/driver/local"
)

const Name = `s3`

var _ upload.Storer = &Filesystem{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, subdir string) (upload.Storer, error) {
		return NewFilesystem(ctx, subdir)
	})
}

func NewFilesystem(ctx context.Context, subdir string) (*Filesystem, error) {
	m, err := model.GetCloudStorage(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	mgr := s3client.New(m.NgingCloudStorage, 0)
	if _, err := mgr.Connect(); err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	return &Filesystem{
		Filesystem: local.NewFilesystem(ctx, subdir, m.Baseurl),
		model:      m,
		mgr:        mgr,
	}, nil
}

// Filesystem 文件系统存储引擎
type Filesystem struct {
	*local.Filesystem
	model *model.CloudStorage
	mgr   *s3manager.S3Manager
}

// Name 引擎名
func (f *Filesystem) Name() string {
	return Name
}

func (f *Filesystem) ErrIsNotExist(err error) bool {
	return f.mgr.ErrIsNotExist(err)
}

// Exists 判断文件是否存在
func (f *Filesystem) Exists(file string) (bool, error) {
	return f.mgr.Exists(context.Background(), file)
}

// FileInfo 获取文件信息
func (f *Filesystem) FileInfo(file string) (os.FileInfo, error) {
	objectInfo, err := f.mgr.Stat(context.Background(), file)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	return s3manager.NewFileInfo(objectInfo), nil
}

// SendFile 下载文件
func (f *Filesystem) SendFile(ctx echo.Context, file string) error {
	ctx.Request().Form().Set(`inline`, `1`)
	err := f.mgr.Download(ctx, file)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// FileDir 物理路径文件夹
func (f *Filesystem) FileDir(subpath string) string {
	return path.Join(uploadLibrary.UploadURLPath, f.Subdir, subpath)
}

// Put 上传文件
func (f *Filesystem) Put(dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = f.FileDir(dstFile)
	//viewURL = `[storage:`+param.AsString(f.model.Id)+`]`+f.URLDir(dstFile)
	viewURL = f.PublicURL(dstFile)
	err = f.mgr.Put(context.Background(), src, savePath, size)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return
}

// Get 获取文件读取接口
func (f *Filesystem) Get(dstFile string) (io.ReadCloser, error) {
	object, err := f.mgr.Get(context.Background(), dstFile)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	exists, err := f.mgr.StatIsExists(object.Stat())
	if !exists {
		err = os.ErrNotExist
	}
	return object, err
}

// Delete 删除文件
func (f *Filesystem) Delete(dstFile string) error {
	err := f.mgr.Remove(context.Background(), dstFile)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// DeleteDir 删除文件夹及其内部文件
func (f *Filesystem) DeleteDir(dstDir string) error {
	err := f.mgr.RemoveDir(context.Background(), dstDir)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// Move 移动文件
func (f *Filesystem) Move(src, dst string) error {
	err := f.mgr.Rename(context.Background(), src, dst)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// Close 关闭连接
func (f *Filesystem) Close() error {
	return nil
}

// FixURL 改写文件网址
func (f *Filesystem) FixURL(content string, embedded ...bool) string {
	rowsByID := f.model.CachedList()
	return uploadLibrary.ReplacePlaceholder(content, func(id string) string {
		r, y := rowsByID[id]
		if !y {
			return ``
		}
		return r.Baseurl
	})
}
