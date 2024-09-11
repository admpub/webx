package image

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/coscms/webcore/cmd/bootconfig"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	modelFile "github.com/coscms/webcore/model/file"
	"github.com/coscms/webcore/registry/upload"
	uploadChecker "github.com/coscms/webcore/registry/upload/checker"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
)

// Proxy http://www.coscms.com/image/proxy/A?image=/public/upload/film/370/72142867778240512.jpg&size=640x351
// {{ImageProxyURL `size`, `640x351`, `image`, `/public/upload/film/370/72142867778240512.jpg` }}
func Proxy(ctx echo.Context) error {
	token := ctx.Param("token")
	if token != uploadChecker.Token(ctx.Queries()) {
		return echo.ErrForbidden
	}
	image := ctx.Query("image")
	size := ctx.Query("size")
	extension := ctx.Query("ex")
	quality := ctx.Queryx(`quality`, `75`).Int()
	if quality < 1 || quality > 100 {
		quality = 75
	}
	sizes := strings.SplitN(size, `x`, 2)
	var width, height float64
	switch len(sizes) {
	case 2:
		width = param.AsFloat64(sizes[0])
		height = param.AsFloat64(sizes[1])
	case 1:
		width = param.AsFloat64(sizes[0])
		height = width
	}
	if width <= 0 || height <= 0 {
		return ctx.NewError(code.InvalidParameter, `Invalid size`).SetZone(`size`)
	}
	var suffix string
	if len(extension) > 0 {
		if strings.Contains(ctx.Header(echo.HeaderAccept), "image/"+extension) {
			suffix = `.` + extension
		}
	}
	mappingKey := fmt.Sprintf(`%s|%vx%v|%s`, image, width, height, suffix)
	if v, ok := mapping.Load(mappingKey); ok {
		if ck, ok := v.(string); ok {
			if v, ok := cached.Load(ck); ok {
				return ctx.ServeCallbackContent(func(ctx echo.Context) (io.Reader, error) {
					return bytes.NewBuffer(v.([]byte)), nil
				}, path.Base(ck), time.Unix(0, 0), bootconfig.HTTPCacheMaxAge)
			}
		}
	}
	// 查询文件记录
	fileM := modelFile.NewFile(ctx)
	err := fileM.Get(nil, db.Cond{`view_url`: image})
	if err != nil {
		if err == db.ErrNoMoreRows {
			return echo.ErrNotFound
		}
		return err
	}
	subdir := fileM.Subdir
	if len(subdir) == 0 {
		subdir = uploadLibrary.ParseSubdir(image)
	}
	storer, err := upload.NewStorer(ctx, subdir)
	if err != nil {
		return err
	}
	defer storer.Close()
	viewURL := image
	image = storer.URLToPath(image)
	thumbURL := tplfunc.AddSuffix(image, fmt.Sprintf(`_%v_%v`, width, height))
	thumbOtherFormatURL := thumbURL + suffix
	thumbM := modelFile.NewThumb(ctx)
	cond := db.And(
		db.Cond{`view_url`: viewURL},
		db.Cond{`width`: width},
		db.Cond{`height`: height},
	)
	err = thumbM.Get(nil, cond)
	if err == nil {
		if len(viewURL) > 7 {
			switch viewURL[0:7] {
			case `http://`, `https:/`:
				return ctx.Redirect(viewURL)
			}
		}
	} else if err != db.ErrNoMoreRows {
		return err
	}
	return ctx.ServeCallbackContent(func(ctx echo.Context) (io.Reader, error) {
		return GetFileReader(ctx, &config{
			mappingKey:          mappingKey,
			thumbM:              thumbM,
			storer:              storer,
			image:               image,
			suffix:              suffix,
			thumbURL:            thumbURL,
			thumbOtherFormatURL: thumbOtherFormatURL,
			viewURL:             viewURL,
			quality:             quality,
			width:               width,
			height:              height,
		})
	}, path.Base(thumbOtherFormatURL), time.Unix(0, 0), bootconfig.HTTPCacheMaxAge)
}
