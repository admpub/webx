package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/listener/upload/friendlink"
	"github.com/admpub/webx/application/model/official"
)

func FriendlinkFormFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(options, formfilter.Exclude(`updated`, `created`, `verifyTime`, `verifyFailCount`, `verifyResult`, `returnTime`, `returnCount`, `host`))
	return formfilter.Build(options...)
}

func FriendlinkIndex(ctx echo.Context) error {
	m := official.NewFriendlink(ctx)
	cond := &db.Compounds{}
	common.SelectPageCond(ctx, cond, `id`, `name%`)
	list, err := m.ListPage(cond, `-id`)
	ctx.Set(`listData`, list)
	return ctx.Render(`official/article/friendlink_index`, handler.Err(ctx, err))
}

func FriendlinkAdd(ctx echo.Context) error {
	var err error
	m := official.NewFriendlink(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonFriendlink, FriendlinkFormFilter())
		if err == nil {
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/article/friendlink_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonFriendlink, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	SetFriendlinkFormData(ctx)
	ctx.Set(`activeURL`, `/official/article/friendlink_index`)
	ctx.Set(`title`, ctx.T(`添加链接`))
	return ctx.Render(`official/article/friendlink_edit`, handler.Err(ctx, err))
}

func FriendlinkEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := official.NewFriendlink(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonFriendlink, FriendlinkFormFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/article/friendlink_index`))
			}
		}
	} else if ctx.IsAjax() {
		process := ctx.Query(`process`)
		if len(process) > 0 {
			m.Process = process
			data := ctx.Data()
			err = m.UpdateField(nil, `process`, process, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonFriendlink, ``, echo.LowerCaseFirstLetter)
	}
	SetFriendlinkFormData(ctx)
	ctx.Set(`activeURL`, `/official/article/friendlink_index`)
	ctx.Set(`title`, ctx.T(`修改链接`))
	return ctx.Render(`official/article/friendlink_edit`, handler.Err(ctx, err))
}

func SetFriendlinkFormData(ctx echo.Context) {
	cateM := official.NewCategory(ctx)
	categoryList := cateM.ListIndent(cateM.ListAllParent(`friendlink`, 0))
	ctx.Set(`categoryList`, categoryList)
	ctx.Set(`thumbSize`, friendlink.FriendlinkLogoThumbnail)
}

func FriendlinkDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := official.NewFriendlink(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/article/friendlink_index`))
}
