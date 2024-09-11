package image

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/nging/v5/application/handler/manager"
	modelFile "github.com/coscms/webcore/model/file"
	"github.com/coscms/webcore/registry/upload/convert"
	"github.com/coscms/webcore/registry/upload/driver"
	"github.com/webx-top/echo"
	"golang.org/x/sync/singleflight"
)

var (
	sg      singleflight.Group
	cached  sync.Map
	mapping sync.Map
)

type config struct {
	thumbM                                                *modelFile.Thumb
	storer                                                driver.Storer
	image, suffix, thumbURL, thumbOtherFormatURL, viewURL string
	quality                                               int
	width, height                                         float64
	mappingKey                                            string
}

func (cfg *config) convertOtherFormatThumb(convert convert.Convert, data []byte) (buf *bytes.Buffer, err error) {
	var val interface{}
	val, err, _ = sg.Do(`convert:`+cfg.thumbOtherFormatURL, func() (interface{}, error) {
		buf, err := convert(io.NopCloser(bytes.NewReader(data)), cfg.quality)
		if err != nil {
			return nil, fmt.Errorf(`%s: %w`, cfg.thumbOtherFormatURL, err)
		}
		b := buf.Bytes()
		cached.Store(cfg.thumbOtherFormatURL, b)
		mapping.Store(cfg.mappingKey, cfg.thumbOtherFormatURL)
		_, _, err = cfg.storer.Put(cfg.storer.URLToFile(cfg.thumbOtherFormatURL), buf, int64(len(b)))
		if err != nil {
			return nil, fmt.Errorf(`%s: %w`, cfg.thumbOtherFormatURL, err)
		}
		return bytes.NewBuffer(b), err
	})
	if err == nil {
		buf = val.(*bytes.Buffer)
	}
	return
}

func (cfg *config) getFileReaderFromStorer(imageURL string) (buf *bytes.Buffer, err error) {
	var val interface{}
	val, err, _ = sg.Do(`get:`+imageURL, func() (interface{}, error) {
		fr, err := cfg.storer.Get(imageURL)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(fr)
		fr.Close()
		if err != nil {
			return nil, fmt.Errorf(`%s: %w`, imageURL, err)
		}
		return bytes.NewBuffer(b), err
	})
	if err == nil {
		buf = val.(*bytes.Buffer)
	}
	return
}

func GetFileReader(ctx echo.Context, cfg *config) (io.Reader, error) {
	var err error

	// ======================
	// 直接获取该文件
	// ======================
	if len(cfg.suffix) > 0 {
		if v, ok := cached.Load(cfg.thumbOtherFormatURL); ok {
			return bytes.NewBuffer(v.([]byte)), err
		}
		buf, err := cfg.getFileReaderFromStorer(cfg.thumbOtherFormatURL)
		if err != nil {
			if !cfg.storer.ErrIsNotExist(err) {
				return nil, fmt.Errorf(`%s: %w`, cfg.thumbOtherFormatURL, err)
			}
		} else {
			cached.Store(cfg.thumbOtherFormatURL, buf.Bytes())
			mapping.Store(cfg.mappingKey, cfg.thumbOtherFormatURL)
			return buf, err
		}
	}

	// ======================
	// 从缩略图获取后转换
	// ======================
	var thumbBytes []byte
	if v, ok := cached.Load(cfg.thumbURL); ok {
		thumbBytes = v.([]byte)
	} else {
		buf, err := cfg.getFileReaderFromStorer(cfg.thumbURL)
		if err != nil {
			if !cfg.storer.ErrIsNotExist(err) {
				return nil, fmt.Errorf(`%s: %w`, cfg.thumbURL, err)
			}
		} else {
			thumbBytes = buf.Bytes()
			cached.Store(cfg.thumbURL, thumbBytes)
		}
	}
	if len(thumbBytes) > 0 {
		if len(cfg.suffix) > 0 {
			if convert, ok := convert.GetConverter(cfg.suffix); ok {
				return cfg.convertOtherFormatThumb(convert, thumbBytes)
			}
		}
		mapping.Store(cfg.mappingKey, cfg.thumbURL)
		return bytes.NewBuffer(thumbBytes), err
	}

	// ======================
	// 从原图获取后生成缩略图后转换
	// ======================
	var imageBytes []byte
	if v, ok := cached.Load(cfg.image); ok {
		imageBytes = v.([]byte)
	} else {
		buf, err := cfg.getFileReaderFromStorer(cfg.image)
		if err != nil {
			if !cfg.storer.ErrIsNotExist(err) {
				return nil, fmt.Errorf(`%s: %w`, cfg.image, err)
			}
		} else {
			imageBytes = buf.Bytes()
			cached.Store(cfg.image, imageBytes)
		}
	}
	if len(imageBytes) > 0 {
		if len(cfg.suffix) > 0 {
			if convert, ok := convert.GetConverter(cfg.suffix); ok {
				return cfg.convertOtherFormatThumb(convert, imageBytes)
			}
		}
	}
	var sgVal interface{}
	var shared bool
	sgVal, err, shared = sg.Do(`makethumb:`+cfg.image, func() (interface{}, error) {
		fileM := modelFile.NewFile(ctx)
		err = fileM.GetByViewURL(cfg.viewURL)
		if err != nil {
			return nil, fmt.Errorf(`%s: %w`, cfg.viewURL, err)
		}
		opt := &imageproxy.Options{
			Width:     cfg.width,  //缩略图宽度
			Height:    cfg.height, //缩略图高度
			Quality:   cfg.quality,
			ScaleUp:   true,
			SmartCrop: true,
		}
		cropOpt := &modelFile.CropOptions{
			Options:          opt,
			File:             fileM.NgingFile,
			SrcReader:        bytes.NewReader(imageBytes),
			Storer:           cfg.storer,
			DestFile:         cfg.storer.URLToFile(cfg.thumbURL),
			FileMD5:          ``,
			WatermarkOptions: manager.GetWatermarkOptions(),
		}
		err = cfg.thumbM.Crop(cropOpt) // 裁剪图片并保存
		if err != nil {
			return nil, err
		}
		byteReader := cropOpt.ThumbData()
		byteReader.Seek(0, 0)
		thumbBytes, err = io.ReadAll(byteReader)
		if err != nil {
			return nil, err
		}
		cached.Store(cfg.thumbURL, thumbBytes)
		byteReader.Seek(0, 0)
		return byteReader, err
	})
	if err != nil {
		return nil, err
	}
	byteReader := sgVal.(*bytes.Reader)
	if len(cfg.suffix) > 0 {
		if convert, ok := convert.GetConverter(cfg.suffix); ok {
			if shared || len(thumbBytes) == 0 {
				thumbBytes, err = io.ReadAll(byteReader)
				if err != nil {
					return nil, err
				}
			}
			return cfg.convertOtherFormatThumb(convert, thumbBytes)
		}
	}
	mapping.Store(cfg.mappingKey, cfg.thumbURL)
	return byteReader, err
}
