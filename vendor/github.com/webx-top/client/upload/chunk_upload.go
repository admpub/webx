package upload

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/admpub/log"
	"github.com/webx-top/com"
)

// 分片上传
func (c *ChunkUpload) Upload(r *http.Request, opts ...ChunkInfoOpter) (int64, error) {
	info := &ChunkInfo{
		FormField: `file`,
	}
	for _, opt := range opts {
		opt(info)
	}
	info.Init(r.FormValue, r.Header.Get)
	if !c.IsSupported(info) {
		return 0, ErrChunkUnsupported
	}
	// 获取上传文件
	upFile, fileHeader, err := r.FormFile(info.FormField)
	if err != nil {
		return 0, fmt.Errorf("上传文件错误: %w", err)
	}
	info.FileName = fileHeader.Filename
	info.CurrentSize = uint64(fileHeader.Size)
	defer upFile.Close()
	return c.ChunkUpload(info, upFile)
}

func (c *ChunkUpload) IsSupported(info ChunkInfor) bool {
	err := c.check(info, true)
	if err == nil {
		return true
	}
	return !errors.Is(err, ErrChunkUnsupported)
}

func (c *ChunkUpload) check(info ChunkInfor, ignoreCurrentSize ...bool) error {
	if info.GetFileTotalBytes() < 1 {
		return fmt.Errorf(`%w: FileTotalBytes less than 1`, ErrChunkUnsupported)
	}
	if (len(ignoreCurrentSize) == 0 || !ignoreCurrentSize[0]) && info.GetCurrentSize() < 1 {
		return fmt.Errorf(`%w: CurrentSize less than 1`, ErrChunkUnsupported)
	}
	if info.GetFileChunkBytes() < 1 {
		return fmt.Errorf(`%w: FileChunkBytes less than 1`, ErrChunkUnsupported)
	}
	if info.GetFileTotalChunks() < 1 {
		return fmt.Errorf(`%w: FileTotalChunks less than 1`, ErrChunkUnsupported)
	}
	return nil
}

func (c *ChunkUpload) ChunkFilename(chunkIndex int) string {
	return filepath.Join(c.TempDir, c.GetUIDString(), fmt.Sprintf("%s_%d", c.fileOriginalName, chunkIndex))
}

// 分片上传
func (c *ChunkUpload) ChunkUpload(info ChunkInfor, upFile io.ReadSeeker) (int64, error) {
	if err := c.check(info); err != nil {
		return 0, err
	}

	c.fileOriginalName = filepath.Base(info.GetFileName())
	if len(c.savePath) > 0 && filepath.Base(c.savePath) == c.fileOriginalName {
		fi, err := os.Stat(c.savePath)
		if err == nil && fi.Size() == int64(info.GetFileTotalBytes()) {
			c.setSaveSize(fi.Size())
			return 0, ErrFileUploadCompleted
		}
	}

	chunkSize := int64(info.GetCurrentSize())

	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)

	if err := os.MkdirAll(chunkFileDir, os.ModePerm); err != nil {
		return 0, err
	}

	// 新文件创建
	filePath := filepath.Join(chunkFileDir, fmt.Sprintf("%s_%d", c.fileOriginalName, info.GetChunkIndex()))
	if log.IsEnabled(log.LevelDebug) {
		log.Debug(filePath+`: `, com.Dump(info, false))
	}

	// 获取现在文件大小
	fi, err := os.Stat(filePath)
	var size int64
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
	} else {
		size = fi.Size()
	}
	// 判断文件是否传输完成
	if size == chunkSize {
		return 0, fmt.Errorf("%w: %s (size: %d bytes)", ErrChunkUploadCompleted, filepath.Base(filePath), size)
	}
	start := size
	saveStart := size
	offset := int64(info.GetChunkOffsetBytes())
	if offset > 0 {
		start = 0
		saveStart = offset
	}

	// 进行断点上传
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("%w: %s: %v", ErrChunkHistoryOpenFailed, filePath, err)
	}

	// 将数据写入文件
	total, err := uploadFile(upFile, start, file, saveStart)

	file.Close()

	if err == nil && total == chunkSize {
		var finished bool
		finished, err = c.isFinish(info, c.fileOriginalName)
		if finished {
			err = c.MergeAll(info.GetFileTotalChunks(), info.GetFileChunkBytes(), info.GetFileTotalBytes(), c.fileOriginalName)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return total, err
}

// 上传文件
func uploadFile(upfile io.ReadSeeker, upSeek int64, file *os.File, fSeek int64) (int64, error) {
	// 设置上传偏移量
	upfile.Seek(upSeek, 0)
	// 设置文件偏移量
	file.Seek(fSeek, 0)
	return WriteTo(upfile, file)
}

func WriteTo(r io.Reader, w io.Writer) (n int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			if err == nil {
				err = fmt.Errorf(`%v`, e)
				return
			}
			err = fmt.Errorf(`%w: %v`, err, e)
		}
	}()
	data := make([]byte, 1024)
	n, err = io.CopyBuffer(w, r, data)
	return
}
