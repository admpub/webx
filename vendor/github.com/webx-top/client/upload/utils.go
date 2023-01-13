package upload

import (
	"io"
	"mime/multipart"
	"path"

	"github.com/webx-top/client/upload/watermark"
	"github.com/webx-top/image"
)

type SaveBeforeHook func(file multipart.File, result *Result, options *Options) (newFile multipart.File, size int64, err error)
type ReadBeforeHook func(result *Result) (err error)

type Options struct {
	ClientName       string
	Result           *Result
	Storer           Storer
	WatermarkOptions *image.WatermarkOptions
	ReadBefore       []ReadBeforeHook
	SaveBefore       []SaveBeforeHook
	Checker          Checker
	Callback         Callback
	MaxSize          int64 // 单位字节 (0 代表未设置，小于 0 代表不限制)
}

type OptionsSetter func(options *Options)

func OptClientName(clientName string) OptionsSetter {
	return func(options *Options) {
		options.ClientName = clientName
	}
}

func OptResult(result *Result) OptionsSetter {
	return func(options *Options) {
		options.Result = result
	}
}

func OptStorer(storer Storer) OptionsSetter {
	return func(options *Options) {
		options.Storer = storer
	}
}

func OptWatermarkOptions(wmOpt *image.WatermarkOptions) OptionsSetter {
	return func(options *Options) {
		options.WatermarkOptions = wmOpt
		options.SaveBefore = append(options.SaveBefore, ImageAddWatermark(wmOpt))
	}
}

func OptReadBefore(hooks ...ReadBeforeHook) OptionsSetter {
	return func(options *Options) {
		options.ReadBefore = append(options.ReadBefore, hooks...)
	}
}

func OptSaveBefore(hooks ...SaveBeforeHook) OptionsSetter {
	return func(options *Options) {
		options.SaveBefore = append(options.SaveBefore, hooks...)
	}
}

func OptChecker(checker Checker) OptionsSetter {
	return func(options *Options) {
		options.Checker = checker
	}
}

func OptCallback(callback Callback) OptionsSetter {
	return func(options *Options) {
		options.Callback = callback
	}
}

func OptMaxSize(maxSize int64) OptionsSetter {
	return func(options *Options) {
		options.MaxSize = maxSize
	}
}

type ReaderAndSizer interface {
	io.Reader
	Sizer
}

func AsFile(body ReadCloserWithSize) (file multipart.File, err error) {
	var oldBody []byte
	oldBody, err = io.ReadAll(body)
	if err != nil {
		return
	}
	file = watermark.Bytes2file(oldBody)
	return
}

// ImageAddWatermark 图片加水印
func ImageAddWatermark(watermarkOptions *image.WatermarkOptions) SaveBeforeHook {
	return func(file multipart.File, result *Result, options *Options) (newFile multipart.File, size int64, err error) {
		if result.FileType.String() != `image` || watermarkOptions == nil || !watermarkOptions.IsEnabled() {
			return file, -1, nil
		}
		b, err := io.ReadAll(file)
		if err != nil {
			return file, -1, err
		}
		b, err = watermark.Bytes(b, path.Ext(result.FileName), watermarkOptions)
		if err != nil {
			return file, -1, err
		}
		file.Seek(0, 0)
		newFile = watermark.Bytes2file(b)
		size = int64(len(b))
		return
	}
}
