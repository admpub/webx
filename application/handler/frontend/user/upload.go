package user

import (
	"fmt"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/handler/manager"
	"github.com/admpub/nging/v5/application/handler/manager/file"
	"github.com/admpub/nging/v5/application/library/config"
	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	"github.com/admpub/nging/v5/application/registry/upload"
	"github.com/admpub/nging/v5/application/registry/upload/checker"
	"github.com/admpub/webx/application/middleware/sessdata"
)

func setUploadURL(ctx echo.Context) error {
	subdir := ctx.Form(`subdir`, `customer`)
	if !upload.AllowedSubdir(subdir) {
		return ctx.NewError(code.InvalidParameter, `无效的subdir值`)
	}
	ctx.Set(`subdir`, subdir)
	ctx.Set(`uploadURL`, checker.FrontendUploadURL(subdir))
	return nil
}

func init() {
	upload.Subdir.Add(`customer`, `用户文件`)
	//manager.SetWatermark(``)
}

func OwnerData(ctx echo.Context) (ownerType string, ownerID uint64) {
	ownerType = `customer`
	customer := sessdata.Customer(ctx)
	if customer == nil {
		user := handler.User(ctx)
		if user != nil {
			ownerType = `user`
			ownerID = uint64(user.Id)
		}
	} else {
		ownerID = customer.Id
	}
	return
}

// Upload 上传文件
func Upload(ctx echo.Context) error {
	ownerType, ownerID := OwnerData(ctx)
	if ownerID < 1 {
		ctx.Data().SetError(ctx.NewError(code.Unauthenticated, `请先登录`))
		return ctx.Redirect(sessdata.URLFor(`/sign_in`))
	}
	uploadCfg := uploadLibrary.Get()
	return manager.UploadByOwner(ctx, ownerType, ownerID, func(result *uploadClient.Result) error { // 自动根据文件类型获取最大上传尺寸
		fileType := result.FileType.String()
		maxSize := uploadCfg.MaxSizeBytes(fileType)
		if maxSize <= 0 {
			maxSize = config.FromFile().GetMaxRequestBodySize()
		}
		if result.FileSize > int64(maxSize) {
			return fmt.Errorf(`%w: %v`, uploadClient.ErrFileTooLarge, com.FormatBytes(maxSize))
		}
		return nil
	})
}

// Crop 图片裁剪
func Crop(ctx echo.Context) error {
	ownerType, ownerID := OwnerData(ctx)
	if ownerID < 1 {
		ctx.Data().SetError(ctx.NewError(code.Unauthenticated, `请先登录`))
		return ctx.Redirect(sessdata.URLFor(`/sign_in`))
	}
	return manager.CropByOwner(ctx, ownerType, ownerID)
}

// Finder 文件选择
func Finder(ctx echo.Context) error {
	if err := setUploadURL(ctx); err != nil {
		return err
	}
	ownerType, ownerID := OwnerData(ctx)
	if ownerID < 1 {
		ctx.Data().SetError(ctx.NewError(code.Unauthenticated, `请先登录`))
		return ctx.Redirect(sessdata.URLFor(`/sign_in`))
	}
	err := file.List(ctx, ownerType, ownerID)
	var suffix string
	if ctx.Formx(`partial`).Bool() {
		suffix = `_partial`
	} else {
		ctx.Set(`subdirList`, upload.Subdir.Slice())
	}
	multiple := ctx.Formx(`multiple`).Bool()
	ctx.Set(`dialog`, true)
	ctx.Set(`multiple`, multiple)
	return ctx.Render(`user/file/list`+suffix, err)
}
