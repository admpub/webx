package user

import (
	"fmt"
	"path"
	"strings"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler/manager"
	"github.com/admpub/nging/v5/application/handler/manager/file"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	"github.com/coscms/webcore/registry/upload"
	"github.com/coscms/webcore/registry/upload/checker"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
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
		user := backend.User(ctx)
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

	var verify []func(result *uploadClient.Result) error

	if ownerType == `customer` {
		customer := sessdata.Customer(ctx)
		m := modelCustomer.NewCustomer(ctx)
		cfg := m.GetUploadConfig(customer)
		verify = append(verify, func(result *uploadClient.Result) error {
			if !cfg.CanUploadSVG {
				extension := path.Ext(result.FileName)
				if strings.EqualFold(extension, `.svg`) {
					return ctx.NewError(code.NonPrivileged, `您没有上传SVG图片的权限`)
				}
			}
			return nil
		})
		if ctx.Form(`subdir`) == `avatar` {
			if !cfg.CanUploadAvatar {
				return ctx.NewError(code.NonPrivileged, `您没有上传头像的权限`)
			}
		} else {
			err := m.Get(func(r db.Result) db.Result {
				return r.Select(`id`, `file_num`, `file_size`)
			}, `id`, customer.Id)
			if err != nil {
				if err == db.ErrNoMoreRows {
					return ctx.NewError(code.UserNotFound, ``)
				}
				return err
			}
			if cfg != nil {
				if m.FileNum+1 > cfg.MaxTotalNum {
					return ctx.NewError(
						code.Failure,
						`上传失败。您的文件数量已满(%s)`,
						cfg.MaxTotalNum,
					)
				}
				verify = append(verify, func(result *uploadClient.Result) error {
					sz := uint64(result.FileSize)
					if sz+m.FileSize > cfg.MaxTotalSizeBytes() {
						return ctx.NewError(
							code.Failure,
							`上传失败。本文件尺寸(%s)加上您已占用空间(%s)超过角色限制(%s)`,
							com.FormatBytes(result.FileSize, 2, true),
							com.FormatBytes(m.FileSize, 2, true),
							cfg.MaxTotalSize,
						)
					}
					m.FileSize += sz
					return nil
				})
			}
		}
	}
	uploadCfg := uploadLibrary.Get()
	return manager.UploadByOwner(ctx, ownerType, ownerID, func(result *uploadClient.Result) (err error) { // 自动根据文件类型获取最大上传尺寸
		fileType := result.FileType.String()
		maxSize := uploadCfg.MaxSizeBytes(fileType)
		if maxSize <= 0 {
			maxSize = config.FromFile().GetMaxRequestBodySize()
		}
		if result.FileSize > int64(maxSize) {
			err = fmt.Errorf(`%w: %v`, uploadClient.ErrFileTooLarge, com.FormatBytes(maxSize))
			return
		}
		for _, fn := range verify {
			if err = fn(result); err != nil {
				return
			}
		}
		return
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
