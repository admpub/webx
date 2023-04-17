package download

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/errors"
	godl "github.com/admpub/go-download/v2"
	"github.com/admpub/log"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/client/upload/watermark"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/nging/v5/application/handler/manager"
	"github.com/admpub/nging/v5/application/library/notice"
	modelFile "github.com/admpub/nging/v5/application/model/file"
	"github.com/admpub/nging/v5/application/registry/upload"
	"github.com/admpub/nging/v5/application/registry/upload/checker"
	"github.com/admpub/nging/v5/application/registry/upload/driver"
)

func Download(ctx echo.Context, options ...Options) (*uploadClient.Result, string, error) {
	options = append(options, OptionsWatermark(manager.GetWatermarkOptions()))
	return AdvanceDownload(ctx, options...)
}

func AdvanceDownload(ctx echo.Context, options ...Options) (*uploadClient.Result, string, error) {
	config := NewConfig(ctx)
	for _, option := range options {
		option(config)
	}
	var thumbURL string
	result := &uploadClient.Result{
		FileType: uploadClient.FileType(config.PrepareData.FileType),
	}
	if len(config.FileURL) == 0 {
		return result, thumbURL, nil
	}
	if config.NoticeSender == nil {
		config.NoticeSender = notice.DefaultNoticer
	}
	fileURL := config.FileURL
	if idx := strings.LastIndex(fileURL, `?`); idx != -1 {
		fileURL = fileURL[:idx]
	} else if idx := strings.LastIndex(fileURL, `#`); idx != -1 {
		fileURL = fileURL[:idx]
	}
	if sendErr := config.NoticeSender(`下载图片 "`+fileURL+`"`, 1, config.Progress); sendErr != nil {
		return result, thumbURL, sendErr
	}
	result.FileName = path.Base(fileURL)
	extension := path.Ext(fileURL)
	godlCfg := &godl.Options{
		Client: func() http.Client {
			return *Client
		},
	}
	if config.MaxMB > 0 {
		godlCfg.MaxSize = config.MaxMB * godl.MB
	}
	dl := func() ([]byte, error) {
		file, err := godl.Open(fileURL, godlCfg)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return io.ReadAll(file)
	}
	b, err := dl()
	if err != nil {
		printRetryMsg := func(i int) {
			retryMsg := fmt.Sprintf(`Will retry in %f seconds. %d retries left.`, config.RetryInterval.Seconds(), config.MaxRetries-i)
			config.NoticeSender(`下载图片 "`+fileURL+`" 失败: `+err.Error()+` (`+retryMsg+`)`, 0, config.Progress)
		}
		printRetryMsg(0)
		for i := 1; i <= config.MaxRetries; i++ {
			time.Sleep(config.RetryInterval)
			iStr := strconv.Itoa(i)
			b, err = dl()
			if err == nil {
				config.NoticeSender(`(重试`+iStr+`)下载图片 "`+fileURL+`" 成功`, 1, config.Progress)
				break
			}
			printRetryMsg(i)
		}
	}
	if err != nil {
		return result, thumbURL, errors.WithMessage(err, fileURL)
	}
	defer func() {
		if e := recover(); e != nil {
			config.NoticeSender(ctx.T(`下载图片 "`+fileURL+`" 出错: %v`, e), 0, config.Progress)
			err = fmt.Errorf(`下载图片 "`+fileURL+`" 出错: %v`, e)
		}
	}()

	if config.Watermark != nil && config.Watermark.IsEnabled() {
		if wmb, err := watermark.Bytes(b, extension, config.Watermark); err != nil {
			if sendErr := config.NoticeSender(`下载图片 "`+fileURL+`", 添加水印失败`, 0, config.Progress); sendErr != nil {
				return result, thumbURL, sendErr
			}
			//return result, thumbURL, err
		} else {
			b = wmb
		}
	}
	//入库
	fileM := config.PrepareData.MakeModel(`user`, 1)
	var name, subdir string
	if config.Checkin {
		subdir, name, err = config.PrepareData.Checkin(ctx)
	} else {
		subdir, name, err = checker.DefaultNoCheck(ctx)
	}
	if err != nil {
		return result, thumbURL, err
	}
	var dstFile string
	dstFile, err = manager.SaveFilename(subdir, name, path.Base(fileURL))
	if err != nil {
		return result, thumbURL, errors.WithMessage(err, fileURL)
	}
	err = config.PrepareData.Checker(result, nil)
	if err != nil {
		return result, thumbURL, errors.WithMessage(err, fileURL)
	}
	var storer driver.Storer
	storer, err = config.PrepareData.Storer()
	if err != nil {
		return result, thumbURL, errors.WithMessage(err, fileURL)
	}
	readCloser := watermark.Bytes2file(b)
	defer readCloser.Close()
	// 保存文件
	result.SavePath, result.FileURL, err = storer.Put(dstFile, readCloser, int64(len(b)))
	if err != nil {
		return result, thumbURL, errors.WithMessage(err, fileURL)
	}
	readCloser.Seek(0, 0)
	result.FileSize = int64(len(b))
	result.Md5 = com.ByteMd5(b)
	fileM.SetByUploadResult(result)

	// 记录到数据库
	err = config.PrepareData.DBSaver(fileM, result, readCloser)
	if err != nil {
		if err := storer.Delete(result.SavePath); err != nil {
			log.Error(result.FileURL, `: `, err)
		}
		return result, thumbURL, err
	}
	readCloser.Seek(0, 0)
	thumbURL, err = Crop(result, fileM, readCloser, storer, config)
	return result, thumbURL, err
}

func Crop(
	result *uploadClient.Result,
	fileM *modelFile.File,
	reader io.Reader,
	storer upload.Storer,
	config *Config,
) (string, error) {
	if config.Crop == nil {
		return ``, nil
	}
	thumbM := modelFile.NewThumb(config.Context)
	/*
		cropOptions := &imageproxy.Options{
			//CropX:          x,   //裁剪X轴起始位置
			//CropY:          y,   //裁剪Y轴起始位置
			//CropWidth:      registryFilm.FilmThumbnail.Width,  //裁剪宽度
			//CropHeight:     registryFilm.FilmThumbnail.Height, //裁剪高度
			Width:     registryFilm.FilmThumbnail.Width,  //缩略图宽度
			Height:    registryFilm.FilmThumbnail.Height, //缩略图高度
			Quality:   75,
			ScaleUp:   true,
			SmartCrop: true,
		}
	*/
	thumbURL := tplfunc.AddSuffix(result.FileURL, fmt.Sprintf(`_%v_%v`, config.Crop.Width, config.Crop.Height))
	cropOpt := &modelFile.CropOptions{
		Options:          config.Crop,
		File:             fileM.NgingFile,
		SrcReader:        reader,
		Storer:           storer,
		DestFile:         storer.URLToFile(thumbURL),
		FileMD5:          ``,
		WatermarkOptions: config.Watermark,
	}
	err := thumbM.Crop(cropOpt)
	return thumbURL, err
}
